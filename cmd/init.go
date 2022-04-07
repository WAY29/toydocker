package cmd

import (
	"github.com/WAY29/toydocker/container"
	cli "github.com/jawher/mow.cli"
)

func CmdInit(cmd *cli.Cmd) {
	var (
		command = cmd.StringArg("COMMAND", "", "command to run")
	)

	cmd.Action = func() {
		container.RunContainerinitProcess(*command)
	}
}
