package server

import (
	"net/http"

	"github.com/bynil/sov2ex/pkg/config"
	"github.com/bynil/sov2ex/pkg/log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type GinRecoverErrorWriter struct{}

func (*GinRecoverErrorWriter) Write(p []byte) (n int, err error) {
	log.Error(string(p))
	return len(p), nil
}

func SetupEngine() *gin.Engine {
	engine := gin.New()
	engine.Use(Ginzap(log.GetLogger(), false))
	engine.Use(gin.RecoveryWithWriter(&GinRecoverErrorWriter{}))

	if config.C.EnableCORS {
		corsConfig := cors.DefaultConfig()
		corsConfig.AllowCredentials = true
		corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Cache-Control", "Pragma", "Authorization")
		corsConfig.AllowOriginFunc = func(origin string) bool {
			return true
		}
		engine.Use(cors.New(corsConfig))
	}

	InitCollection()
	registerRouters(engine)
	return engine
}

func registerRouters(engine *gin.Engine) {
	engine.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	engine.GET("/api/search", searchHandler)
}