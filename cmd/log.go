package cmd

import (
	"github.com/WAY29/toydocker/container"
	cli "github.com/jawher/mow.cli"
)

func CmdLogs(cmd *cli.Cmd) {
	var (
		containerID = cmd.StringArg("CONTAINER", "", "Container name or id")
	)

	cmd.Spec = "CONTAINER"

	cmd.Before = container.InitDatabase

	cmd.Action = func() {
		container.ShowContainerLogs(*containerID)
	}

}
