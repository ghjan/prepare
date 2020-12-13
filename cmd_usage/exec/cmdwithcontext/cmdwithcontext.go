package cmdwithcontext

import (
	"context"
	"os/exec"
)

type Result struct {
	Err    error
	Output []byte
}

var (
	cmd *exec.Cmd
)

func UseExecContext(ctx context.Context, cmdString []string, resultChan chan *Result) {
	var (
		output []byte
		err    error
	)
	cmd = exec.CommandContext(ctx, cmdString[0], cmdString[1:]...)
	//执行了命令，捕获了子进程的输出(pipe)
	output, err = cmd.CombinedOutput()
	//把结果传出去
	resultChan <- &Result{
		Err:    err,
		Output: output,
	}
}
