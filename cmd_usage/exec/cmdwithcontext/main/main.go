package main

import (
	"context"
	"fmt"
	"prepare/cmd_usage/exec/cmdwithcontext"
	"time"
)

var (
	ctx        context.Context
	cancelFunc context.CancelFunc
	resultChan chan *cmdwithcontext.Result
)

func main() {
	cmdStrings := [][]string{
		[]string{"bash.exe", "-c", "ls -la"},
		[]string{"bash.exe", "-c", "sleep 1;ls -la"},
		[]string{"bash.exe", "-c", "sleep 2;ls -la"},
	}
	resultChan = make(chan *cmdwithcontext.Result, 100)
	for index, cmdString := range cmdStrings {
		fmt.Printf("index:%d\n", index)
		var res *cmdwithcontext.Result
		ctx, cancelFunc = context.WithCancel(context.TODO())
		go cmdwithcontext.UseExecContext(ctx, cmdString, resultChan)
		time.Sleep(1 * time.Second)
		//取消上下文
		cancelFunc()

		//获取执行结果
		res = <-resultChan
		//打印执行结果
		if res != nil {
			if res.Err != nil {
				fmt.Printf("index:%d,error:%s\n", index, res.Err.Error())
			} else {
				//打印子进程的输出
				fmt.Printf("index:%d,output:%s\n", index, string(res.Output))
			}
		} else {
			fmt.Printf("index:%d,oops!res is nil!\n", index)
		}

	}
}
