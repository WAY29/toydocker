package cmd

import (
	"github.com/WAY29/toydocker/container"
	cli "github.com/jawher/mow.cli"
)

func CmdRun(cmd *cli.Cmd) {
	var (
		tty = cmd.BoolOpt("t tty", false, "enable tty")
	)

	var (
		command = cmd.StringArg("COMMAND", "", "command to run")
	)

	cmd.Spec = "[-t | --tty] COMMAND"

	cmd.Action = func() {
		container.Run(*tty, *command)
	}
}
