package mongodb_conn

import (
	"math/rand"
	"strconv"
	"time"
)

type Student struct {
	Name string
	Age  int
}
type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}

type LogRecord struct {
	JobName   string    `bson:"jobName"`
	Command   string    `bson:"command"`
	Err       string    `bson:"err"`
	Content   string    `bson:"content"`
	TimePoint TimePoint `bson:"timePoint"`
}
type FindByJobName struct {
	JobName string `bson:"jobName"`
}

//startTime小于某时间
//{"$lt":timestamp}
type TimeBeforeCond struct {
	Before int64 `bson:"$lt"`
}

//{"timePoint.startTime":{"$lt":timestamp}}
type DeleteCond struct {
	beforeCond TimeBeforeCond `bson:"timePoint.startTime"`
}

func GenerateLogRecord(jobName string) (record LogRecord) {
	rand.Seed(time.Now().Unix())
	number := strconv.Itoa(rand.Intn(1000000))
	if jobName == "" {
		jobName = "job" + number
	}
	command := "echo hello" + number
	content := "hello" + number
	record = LogRecord{JobName: jobName, Command: command, Content: content,
		TimePoint: TimePoint{StartTime: time.Now().Unix(),
			EndTime: time.Now().Unix() + 10}}
	return
}
