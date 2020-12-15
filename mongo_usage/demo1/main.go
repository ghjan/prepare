package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"prepare/mongo_usage/mongodb_conn"
)

func main() {
	var (
		client     *mongo.Client
		err        error
		database   *mongo.Database
		collection *mongo.Collection
	)
	//1.建立连接
	if client, err = mongodb_conn.GetMongoClient(); err != nil {
		fmt.Printf("mongodb_conn.GetMongoClient error:%s\n", err.Error())
		return
	}

	//2.选择数据库my_db
	database = client.Database("my_db")

	//3.选择表my_collection
	collection = database.Collection("my_collection")
	collection = collection

}
