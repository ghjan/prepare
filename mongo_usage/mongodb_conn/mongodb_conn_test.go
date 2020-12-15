package mongodb_conn

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
	"time"
)

var (
	collection     *mongo.Collection
	dbName         string = "cron"
	collectionName string = "log"
	Job1Name       string = "job1"
)

func TestMongoInserOne(t *testing.T) {
	s1 := GenerateLogRecord(Job1Name)

	insertResult, err := collection.InsertOne(context.TODO(), s1)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	t.Log("Inserted a single document: ", insertResult.InsertedID)

}

func TestMongoInsertMany(t *testing.T) {
	s2 := GenerateLogRecord("")
	s3 := GenerateLogRecord("")
	jobs := []interface{}{s2, s3}
	insertManyResult, err := collection.InsertMany(context.TODO(), jobs)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	t.Log("Inserted multiple documents: ", insertManyResult.InsertedIDs)
}
func TestMongoFind(t *testing.T) {

	collection := mongoClient.Database(dbName).Collection(collectionName)
	// 创建一个LogRecord变量用来接收查询的结果
	var result LogRecord
	//filter := bson.D{{"jobName", Job1Name}}
	filter := &FindByJobName{JobName: Job1Name}
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	t.Logf("Found a single document: %+v\n", result)

	// 查询多个
	// 将选项传递给Find()
	findOptions := options.Find()
	//findOptions.SetSkip(0)  //开始
	//findOptions.SetLimit(2) //返回记录数量

	// 定义一个切片用来存储查询结果
	var results []*LogRecord

	// 把bson.D{{}}作为一个filter来匹配所有文档
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions.SetSkip(0), findOptions.SetLimit(2))
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	// 延迟释放后关闭游标
	defer cur.Close(context.TODO())

	// 查找多个文档返回一个光标
	// 遍历游标允许我们一次解码一个文档
	for cur.Next(context.TODO()) {
		// 创建一个值，将单个文档解码为该值
		var elem LogRecord
		err := cur.Decode(&elem)
		if err != nil {
			t.Error(err)
			t.Fail()
			return
		}
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Logf("Found multiple documents (array of pointers):\n")
	for _, r := range results {
		t.Logf("%#v\n", *r)
	}

	t.Log("Connection to MongoDB closed.")
}

func TestDelete(t *testing.T) {

	var (
		delCond   *DeleteCond
		delResult *mongo.DeleteResult
	)
	// 删除名字是job1的那个
	filter := &FindByJobName{JobName: Job1Name}
	//bson.D{{"jobName", Job1Name}}
	deleteResult1, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Logf("Deleted %v documents in the collection\n", deleteResult1.DeletedCount)
	//删除开始时间早于当前时间的所有日志
	//delete({"timePoint.startTime":{"$lt":当前时间}})

	delCond = &DeleteCond{beforeCond: TimeBeforeCond{Before: time.Now().Unix()}}
	if delResult, err = collection.DeleteMany(context.TODO(), delCond); err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	assert.Equal(t, delResult.DeletedCount > 0, true)
	t.Logf("Deleted %v documents in the collection\n", delResult.DeletedCount)

	// 删除所有
	deleteResult2, err := collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	assert.Equal(t, int64(0), deleteResult2.DeletedCount)
	//t.Logf("Deleted %v documents in the collection\n", deleteResult2.DeletedCount)

}

func TestMain(m *testing.M) {
	fmt.Println("mongodb_conn test begin")
	mongoClient, err := GetMongoClient()
	if err != nil {
		log.Println(err)
	}
	collection = mongoClient.Database(dbName).Collection(collectionName)

	m.Run()
	// 断开连接
	err = mongoClient.Disconnect(context.TODO())
	if err != nil {
		log.Println(err)
	}
	log.Println("Connection to MongoDB closed.")
	fmt.Println("mongodb_conn test end")
}
