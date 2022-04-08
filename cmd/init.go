package cmd

import (
	"os"

	"github.com/WAY29/toydocker/container"
	cli "github.com/jawher/mow.cli"
	log "github.com/sirupsen/logrus"
)

func CmdInit(cmd *cli.Cmd) {
	cmd.Action = func() {
		if os.Args[0] != "/proc/self/exe" {
			log.Error("Can't call this outside")
			cli.Exit(1)
		}

		container.RunContainerInitProcess()
	}
}
