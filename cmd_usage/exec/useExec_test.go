package exec

import (
	"fmt"
	"testing"
)

func TestUseExec(t *testing.T) {
	cmdStrings := [][]string{
		[]string{"cmd", "/C", "dir", "D:\\"},
		[]string{"cmd", "/C", "cd", "D:\\"},
		[]string{"bash.exe", "-c", "sleep 1;ls -la"},
	}
	for index, cmdString := range cmdStrings {
		t.Run(fmt.Sprint(cmdString), func(t *testing.T) {
			t.Logf("index:%d\n", index)
			output, err := UseExec(cmdString)
			if err != nil {
				t.Errorf("error:%s", err.Error())
				t.Fail()
			} else {
				//打印子进程的输出
				t.Log(string(output))
			}
		})
	}
}
