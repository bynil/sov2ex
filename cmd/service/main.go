package main

import (
	"fmt"
	"net/http"

	"github.com/bynil/sov2ex/pkg/config"
	_ "github.com/bynil/sov2ex/pkg/config"
	"github.com/bynil/sov2ex/pkg/es"
	"github.com/bynil/sov2ex/pkg/log"
	"github.com/bynil/sov2ex/pkg/mongodb"
	"github.com/bynil/sov2ex/pkg/server"
)

func main() {
	log.InitLog()
	defer log.Sync()
	log.Info("sov2ex service starting...")
	log.Infof("log level %v", log.Level.String())

	es.InitClient(config.C.ESURL)
	mongodb.InitClient(fmt.Sprintf("mongodb://%v:%v@%v:%v/%v",
		config.C.MongoUser, config.C.MongoPass,
		config.C.MongoHost, config.C.MongoPort,config.C.MongoDBName,
	))

	addr := fmt.Sprintf("%v:%v", config.C.Host, config.C.Port)
	log.Infof("listen address %v", addr)
	engine := server.SetupEngine()
	defer log.Info("engine exit")

	if err := engine.Run(addr); err != nil && err != http.ErrServerClosed {
		log.Error(err)
	}
}
