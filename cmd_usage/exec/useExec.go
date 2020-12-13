package exec

import (
	"os/exec"
)

var (
	cmd *exec.Cmd
)

func UseExec(cmdString []string) (output []byte, err error) {
	//生成cmd
	cmd = exec.Command(cmdString[0],  cmdString[1:]...)
	//执行了命令，捕获了子进程的输出(pipe)
	output, err = cmd.CombinedOutput()
	return
}
