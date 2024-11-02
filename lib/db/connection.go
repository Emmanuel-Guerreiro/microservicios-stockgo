package db

import (
	"context"
	"emmanuel-guerreiro/stockgo/lib"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var client *mongo.Client

func ConnectDatabase() {
	if client != nil {
		return
	}
	clientOpts := options.Client().ApplyURI(lib.GetEnv().MongoURL)
	c, err := mongo.Connect(clientOpts)
	if err != nil {
		log.Panic(err)
	}
	client = c
}

func Get() *mongo.Database {
	if client == nil {
		ConnectDatabase()
	}
	return client.Database("stock")
}

func DisconnectDatabase() {
	if client == nil {
		return
	}
	client.Disconnect(context.TODO())
}

func IsUniqueKeyError(err error) bool {
	if wErr, ok := err.(mongo.WriteException); ok {
		for i := 0; i < len(wErr.WriteErrors); i++ {
			if wErr.WriteErrors[i].Code == 11000 {
				return true
			}
		}
	}
	return false
}
