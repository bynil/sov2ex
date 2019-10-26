package mongodb

import (
	"context"
	"time"

	"github.com/bynil/sov2ex/pkg/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func InitClient(uri string) {
	var err error
	Client, err = mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = Client.Connect(ctx)
	if err != nil {
		log.Panic(err)
	}
}
