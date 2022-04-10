package cmd

import (
	"github.com/WAY29/toydocker/container"
	cli "github.com/jawher/mow.cli"
)

func CmdPS(cmd *cli.Cmd) {
	cmd.Before = container.InitDatabase

	cmd.Action = func() {
		container.PS()
	}
}
