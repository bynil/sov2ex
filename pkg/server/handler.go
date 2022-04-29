package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bynil/sov2ex/pkg/config"
	"github.com/bynil/sov2ex/pkg/es"
	"github.com/bynil/sov2ex/pkg/log"
	"github.com/bynil/sov2ex/pkg/mongodb"
	"github.com/bynil/sov2ex/pkg/utils/int64set"
	"github.com/bynil/sov2ex/pkg/utils/stringset"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/time/rate"
	"gopkg.in/olivere/elastic.v5"
)

const (
	NodeCollectionName = "node"
	TopicAliasName     = "topic"
	TopicTypeName      = "topic"

	SizeDefault    = 10
	SizeMax        = 50
	PagingDepthMax = 1000

	ClauseCountMax   = 30
	KeywordLengthMax = 100

	SortTypeSumup   = "sumup"
	SortTypeCreated = "created"

	OrderTypeDesc = 0
	OrderTypeAsc  = 1

	OperatorTypeOr  = "or"
	OperatorTypeAnd = "and"

	ExcludePrefix          = "-"
	MultipleItemsSeparator = ","

	V2EXUserHomepageFormat = "https://www.v2ex.com/member/%v"

	LimiterWaitTimeMax = 5 * time.Second
)

var (
	c          = cache.New(time.Hour, time.Hour) // for user's searchable status
	decoder    = schema.NewDecoder()
	httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
	limiter = rate.NewLimiter(2, 4)

	ErrUserNotFound       = errors.New("V2EX user not found")
	ErrGetUserInfoFailed  = errors.New("get user info failed")
	ErrRequestLimitExceed = errors.New("exceed the request limit of getting user info")

	SortTypeChoices     = stringset.NewSet(SortTypeSumup, SortTypeCreated)
	OrderTypeChoices    = int64set.NewSet(OrderTypeDesc, OrderTypeAsc)
	OperatorTypeChoices = stringset.NewSet(OperatorTypeOr, OperatorTypeAnd)
)

func init() {
	decoder.IgnoreUnknownKeys(true)
}

type SearchParams struct {
	Keyword  string `schema:"q"`
	From     int64  `schema:"from"`
	Size     int64  `schema:"size"`
	Sort     string `schema:"sort"`
	Order    int64  `schema:"order"`
	Gte      int64  `schema:"gte"`
	Lte      int64  `schema:"lte"`
	Node     string `schema:"node"` // should be split and replaced by node id（int64)
	Operator string `schema:"operator"`
	Username string `schema:"username"`
}

var searchHandler = func(c *gin.Context) {
	params := NewDefaultParams()
	err := decoder.Decode(&params, c.Request.URL.Query())
	if err != nil {
		ReqErrorWithErr(c, http.StatusBadRequest, err)
		return
	}
	if err = validateParams(params); err != nil {
		ReqErrorWithErr(c, http.StatusBadRequest, err)
		return
	}
	rp, err := GenerateRenderParams(params)
	if err != nil {
		if err == ErrUserNotFound {
			ReqErrorWithErr(c, http.StatusNotFound, err)
		} else if err == ErrRequestLimitExceed {
			ReqErrorWithErr(c, http.StatusTooManyRequests, err)
		} else {
			ReqErrorWithErr(c, http.StatusInternalServerError, err)
		}
		return
	}
	var queryBody string
	switch rp.Sort {
	case SortTypeSumup:
		queryBody = RenderScoreSearchBody(rp)
	case SortTypeCreated:
		queryBody = RenderTimeOrderSearchBody(rp)
	default:
		queryBody = RenderScoreSearchBody(rp)
	}

	sr, err := searchInES(queryBody)
	if err != nil {
		log.Error(err)
		ReqErrorWithMessage(c, http.StatusServiceUnavailable, "Elasticsearch error")
		return
	}
	c.JSON(http.StatusOK, sr)
}

func ReqErrorWithMessage(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, map[string]interface{}{"message": msg})
}

func ReqErrorWithErr(c *gin.Context, code int, err error) {
	ReqErrorWithMessage(c, code, err.Error())
}

func NewDefaultParams() SearchParams {
	return SearchParams{
		Keyword:  "",
		From:     0,
		Size:     SizeDefault,
		Sort:     SortTypeSumup,
		Order:    OrderTypeDesc,
		Gte:      0,
		Lte:      0,
		Node:     "",
		Operator: OperatorTypeOr,
	}
}

func validateParams(sp SearchParams) (err error) {
	if sp.Keyword == "" {
		return errors.New("missing keyword")
	}
	if len([]rune(sp.Keyword)) > KeywordLengthMax {
		return errors.New("too long keyword")
	}
	if !SortTypeChoices.Contains(sp.Sort) {
		return errors.New("invalid sort")
	}
	if !OrderTypeChoices.Contains(sp.Order) {
		return errors.New("invalid order")
	}
	if !OperatorTypeChoices.Contains(sp.Operator) {
		return errors.New("invalid operator")
	}
	if sp.From < 0 {
		return errors.New("invalid from")
	}
	if sp.Size < 0 {
		return errors.New("invalid size")
	}
	if sp.From+sp.Size > PagingDepthMax {
		return errors.New("too deep paging")
	}
	if sp.Size > SizeMax {
		return errors.New("too large size")
	}
	num, err := analyzeTokenNum(sp.Keyword)
	if err != nil {
		log.Error(err)
		return errors.New("keyword analyzed failed")
	} else if num > ClauseCountMax {
		return errors.Errorf("too long keyword: %v clauses", num)
	}
	return
}

func GenerateRenderParams(sp SearchParams) (rp RenderParams, err error) {
	defer func() {
		log.Debugf("generated render params: %v", rp)
	}()
	rp.SearchParams = sp

	if sp.Node != "" {
		nodes, exclude := parseItemsInParam(sp.Node)
		nodeIds := int64set.NewSet()
		for _, node := range nodes {
			nodeId, nodeErr := findNodeId(node)
			// ignore node error
			if nodeErr == nil {
				nodeIds.Add(nodeId)
			}
		}
		if nodeIds.Length() > 0 {
			if exclude {
				rp.ExcludedNodeIds = nodeIds.GetSlice()
			} else {
				rp.NodeIds = nodeIds.GetSlice()
			}
		}
	}

	if rp.Username != "" && !config.C.DisableUserCheck {
		var info *userInfo
		info, err = getUserInfo(rp.Username)
		log.Infof("try to get userinfo %v, result: %v, err: %v", rp.Username, info, err)
		if err != nil {
			return
		}
		if !info.Found {
			err = ErrUserNotFound
			return
		}
		if info.Searchable {
			rp.Username = info.RealUserName
		} else {
			rp.Lte = 1 // for empty result
		}
	}
	return
}

func parseItemsInParam(param string) (items []string, excluded bool) {
	splits := strings.Split(param, MultipleItemsSeparator)
	for _, split := range splits {
		trim := strings.TrimSpace(split)
		if strings.HasPrefix(trim, ExcludePrefix) {
			excluded = true
			trim = strings.TrimPrefix(trim, ExcludePrefix)
		}
		items = append(items, trim)
	}
	return
}

/*----- MongoDB -----*/
var (
	nodeCollection *mongo.Collection
)

type nodeDoc struct {
	Id int64 `bson:"id"`
}

func InitCollection() {
	nodeCollection = mongodb.Client.Database(config.C.MongoDBName).Collection(NodeCollectionName)
}

// findNodeId search node id in mongodb, node could be node's name, title, title_alternative,
// return error if node not found.
func findNodeId(node string) (nodeId int64, err error) {
	node = strings.TrimSpace(node)
	if node == "" {
		return nodeId, errors.New("empty node name")
	}
	var doc nodeDoc
	filter := bson.M{
		"$or": []map[string]string{
			{"name": node},
			{"title": node},
			{"title_alternative": node},
		}}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = nodeCollection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return
	}
	return doc.Id, nil
}

/*---------------*/

/*----- Elasticsearch -----*/
type SearchResult struct {
	TookInMillis int64        `json:"took"`      // search time in milliseconds
	TotalHits    int64        `json:"total"`     // total number of hits found
	Hits         []*SearchHit `json:"hits"`      // the actual search hits
	TimedOut     bool         `json:"timed_out"` // true if the search timed out
}

type SearchHit struct {
	Score     *float64                   `json:"_score"`    // computed score
	Index     string                     `json:"_index"`    // index name
	Type      string                     `json:"_type"`     // type meta field
	Id        string                     `json:"_id"`       // external or internal
	Sort      []interface{}              `json:"sort"`      // sort information
	Highlight elastic.SearchHitHighlight `json:"highlight"` // highlighter information
	Source    *json.RawMessage           `json:"_source"`   // stored document source
}

func searchInES(query string) (sr *SearchResult, err error) {
	esResult, err := es.Client.Search().Index(TopicAliasName).Type(TopicTypeName).Source(query).Do(context.Background())
	if err != nil {
		return
	}
	sr = &SearchResult{
		TookInMillis: esResult.TookInMillis,
		TotalHits:    esResult.Hits.TotalHits,
		Hits:         make([]*SearchHit, 0),
		TimedOut:     esResult.TimedOut,
	}
	if esResult.Hits != nil && len(esResult.Hits.Hits) > 0 {
		for _, esHit := range esResult.Hits.Hits {
			sh := SearchHit{
				Score:     esHit.Score,
				Index:     esHit.Index,
				Type:      esHit.Type,
				Id:        esHit.Id,
				Sort:      esHit.Sort,
				Highlight: esHit.Highlight,
				Source:    esHit.Source,
			}
			sr.Hits = append(sr.Hits, &sh)
		}
	}
	return
}

func analyzeTokenNum(keyword string) (tokenNum int, err error) {
	resp, err := es.Client.IndexAnalyze().
		Index(TopicAliasName).Text(keyword).
		Analyzer("ik_smart").
		Do(context.Background())
	if err != nil {
		return
	}
	tokenNum = len(resp.Tokens)
	return
}

/*---------------*/

/*----- V2EX -----*/
type userInfo struct {
	RealUserName string
	Searchable   bool
	Found        bool
}

func crawlUserInfo(username string) (info *userInfo, err error) {
	info = new(userInfo)
	if username == "" {
		return
	}

	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(LimiterWaitTimeMax))
	err = limiter.Wait(ctx)
	if err != nil {
		return nil, ErrRequestLimitExceed
	}

	link := fmt.Sprintf(V2EXUserHomepageFormat, username)
	resp, err := httpClient.Get(link)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return
		} else {
			log.Errorf("fetch user %v homepage error, status code is abnormal: %v", username, resp.StatusCode)
			err = ErrGetUserInfoFailed
		}
		return
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return
	}
	notice := doc.Find("td.topic_content").Text()
	info.Searchable = !(notice != "" && strings.Contains(notice, "根据"))
	info.RealUserName = doc.Find("h1").First().Text()
	info.Found = true
	return
}

func getUserInfo(username string) (info *userInfo, err error) {
	usernameLowerCase := strings.TrimSpace(strings.ToLower(username)) // lower-case as key
	if userInfoI, found := c.Get(usernameLowerCase); found {
		return userInfoI.(*userInfo), nil
	}
	info, err = crawlUserInfo(usernameLowerCase)
	if err != nil {
		return
	}
	// not found user will be cached
	c.Set(usernameLowerCase, info, cache.DefaultExpiration)
	return
}

/*---------------*/
