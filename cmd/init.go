package cmd

import (
	"os"

	"github.com/WAY29/toydocker/container"
	cli "github.com/jawher/mow.cli"
	"github.com/sirupsen/logrus"
)

func CmdInit(cmd *cli.Cmd) {
	cmd.Action = func() {
		if os.Args[0] != "/proc/self/exe" {
			logrus.Error("Can't call this outside")
			cli.Exit(1)
		}

		container.RunContainerInitProcess()
	}
}
