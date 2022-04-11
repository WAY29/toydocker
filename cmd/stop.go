package cmd

import (
	"github.com/WAY29/toydocker/container"
	cli "github.com/jawher/mow.cli"
)

func CmdStop(cmd *cli.Cmd) {
	var ()

	var (
		containersID = cmd.StringsArg("CONTAINER", []string{}, "Container name or id")
	)

	cmd.Spec = "CONTAINER..."

	cmd.Before = container.InitDatabase

	cmd.Action = func() {
		container.StopContainer(*containersID)
	}

}
