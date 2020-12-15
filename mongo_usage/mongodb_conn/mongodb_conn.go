package mongodb_conn

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var mongoClient *mongo.Client

func GetMongoClient() (*mongo.Client, error) {
	if mongoClient != nil {
		return mongoClient, nil
	} else {
		return initMongoClient()
	}
}

func initMongoClient() (client *mongo.Client, err error) {
	if mongoClient == nil {
		// Set mongoClient options
		clientOptions1 := options.Client().ApplyURI("mongodb://localhost:27017")
		clientOptions2 := options.Client().SetConnectTimeout(5 * time.Second)
		// Connect to MongoDB
		if client, err = mongo.Connect(context.TODO(), clientOptions1, clientOptions2); err != nil {
			log.Fatalf("initMongoClient, mongo.Connect error:%s\n", err.Error())
			return
		}

		// Check the connection
		if err = client.Ping(context.TODO(), nil); err != nil {
			log.Fatalf("initMongoClient, mongoClient.Ping error:%s\n", err.Error())
			return
		}

		mongoClient = client
		fmt.Println("Connected to MongoDB!")
		return
	}
	return

}
