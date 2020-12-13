package cmdwithcontext

import (
	"context"
	"fmt"
	"testing"
	"time"
)

var (
	ctx        context.Context
	cancelFunc context.CancelFunc
	resultChan chan *Result
)

func TestUseExecContext(t *testing.T) {
	cmdStrings := [][]string{
		[]string{"bash.exe", "-c", "ls -la"},
		[]string{"bash.exe", "-c", "sleep 1;ls -la"},
		[]string{"bash.exe", "-c", "sleep 2;ls -la"},
	}
	resultChan = make(chan *Result, 100)
	for index, cmdString := range cmdStrings {
		t.Run(fmt.Sprint(cmdString), func(t *testing.T) {
			t.Logf("index:%d", index)
			var res *Result
			ctx, cancelFunc = context.WithCancel(context.TODO())
			go UseExecContext(ctx, cmdString, resultChan)
			time.Sleep(1 * time.Second)
			//取消上下文
			cancelFunc()

			//获取执行结果
			res = <-resultChan
			//打印执行结果
			if res != nil {
				if res.Err != nil {
					t.Logf("index:%d,error:%s", index, res.Err.Error())
				} else {
					//打印子进程的输出
					t.Logf("index:%d,output:%s", index, string(res.Output))
				}
			} else {
				t.Logf("index:%d,oops!res is nil!", index)
			}
		})
	}
}
