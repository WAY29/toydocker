package cmd

import (
	"github.com/WAY29/toydocker/container"
	cli "github.com/jawher/mow.cli"
)

func CmdExport(cmd *cli.Cmd) {
	var (
		output = cmd.StringOpt("o output", "", "Write to a file path")
	)

	var (
		containerID = cmd.StringArg("CONTAINER", "", "Container name or id")
	)

	cmd.Spec = "(-o=<output filepath> | --output=<output filepath>) CONTAINER"

	cmd.Before = container.InitDatabase

	cmd.Action = func() {
		container.ExportContainer(*output, *containerID)
	}

}
