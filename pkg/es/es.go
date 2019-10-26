package es

import (
	"github.com/bynil/sov2ex/pkg/log"
	"github.com/pkg/errors"
	"gopkg.in/olivere/elastic.v5"
)

var (
	Client *elastic.Client
)
func InitClient(esURL string) {
	var err error
	Client, err = elastic.NewClient(elastic.SetURL(esURL))
	if err != nil {
		log.Panic(errors.Wrap(err, "init es client error"))
	}

	esVersion, err := Client.ElasticsearchVersion(esURL)
	if err != nil {
		log.Panic(errors.Wrap(err, "es client ping error"))
	}
	log.Infof("Elasticsearch version %s", esVersion)
}
