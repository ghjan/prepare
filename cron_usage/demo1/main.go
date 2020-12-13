package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

func main() {
	var (
		expr     *cronexpr.Expression
		err      error
		now      time.Time
		nextTime time.Time
		finished chan bool
	)
	//哪一分钟（0-59） 哪小时
	//每分钟执行一次
	//if expr, err = cronexpr.Parse("* * * * *"); err != nil {
	//	fmt.Printf("error:%s\n", err.Error())
	//}

	//linux crontab 5个单位 最小只能到分钟粒度 分钟 小时 天 月 星期几 年
	//从0开始 每隔5秒钟 cronexpr可以到秒粒度 年配置(2018-2099)
	if expr, err = cronexpr.Parse("*/5 * * * * * *"); err != nil {
		fmt.Printf("error:%s\n", err.Error())
	}

	//0 5 10
	//当前时间
	now = time.Now()
	//下次调度时间
	nextTime = expr.Next(now)
	fmt.Printf("now:%v\n", now)
	finished=make(chan bool,0)
	time.AfterFunc(nextTime.Sub(now), func() {
		finished <- true
		fmt.Printf("被调度了nextTime:%v\n", nextTime)
	})
	<-finished
	time.Sleep(1 * time.Second)

}
