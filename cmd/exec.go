package cmd

import (
	"os"

	"github.com/WAY29/toydocker/container"
	_ "github.com/WAY29/toydocker/nsenter"
	cli "github.com/jawher/mow.cli"
)

func CmdExec(cmd *cli.Cmd) {
	var ()

	var (
		containerID = cmd.StringArg("CONTAINER", "", "Container name or id")
		command     = cmd.StringsArg("COMMAND", []string{}, "Command to run")
	)

	cmd.Spec = "CONTAINER COMMAND..."

	cmd.Before = container.InitDatabase

	cmd.Action = func() {
		// 如果已经设置了环境变量，证明已经执行了程序，直接返回
		if os.Getenv(container.ENV_EXEC_CMD) != "" || os.Getenv(container.ENV_EXEC_PID) != "" {
			return
		}
		container.ExecContainer(*containerID, *command)
	}

}
