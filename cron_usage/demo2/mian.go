package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

//调度多个crontab任务
type CronJob struct {
	expr     *cronexpr.Expression
	nextTime time.Time //expr.Next(now)
}

func main() {
	//需要有一个调度goroutine，它定时检查所有的cron任务 谁过期了就执行谁
	var (
		cronJob       *CronJob
		expr          *cronexpr.Expression
		now           time.Time
		scheduleTable map[string]*CronJob //调度表 key：任务名字
	)
	scheduleTable = make(map[string]*CronJob, 0)

	//当前时间
	now = time.Now()

	//1.定义两个cronjob
	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}
	//任务保存到调度表里
	scheduleTable["job1"] = cronJob

	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}
	scheduleTable["job2"] = cronJob

	//启动一个调度goroutine
	//检查是否有任务到期，依据nexttime
	go func() {
		var (
			jobName string
			cronJob *CronJob
		)
		//定时检查一下任务调度表
		for {
			now = time.Now()
			for jobName, cronJob = range scheduleTable {
				//判断是否过期
				if cronJob.nextTime.Before(now) || cronJob.nextTime.Equal(now) {
					//启动一个goroutine 执行这个任务
					go func(jobName string) {
						fmt.Println("执行:", jobName)
					}(jobName)
					//计算下一次调度时间
					cronJob.nextTime = cronJob.expr.Next(now)
					fmt.Println(jobName, "下次执行时间:", cronJob.nextTime)
				}
			} //inner for
			//睡眠100ms
			//time.Sleep(100*time.Millisecond)
			select {
			case <-time.NewTimer(100 * time.Millisecond).C: //将在100ms可读，返回
			}
		} //outer for

	}()
	time.Sleep(100 * time.Second)
}
