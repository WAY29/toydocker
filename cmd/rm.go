package cmd

import (
	"github.com/WAY29/toydocker/container"
	cli "github.com/jawher/mow.cli"
)

func CmdRemove(cmd *cli.Cmd) {
	var (
		force = cmd.BoolOpt("f force", false, "Force the removal of a running container (uses SIGKILL)")
	)

	var (
		containerID = cmd.StringsArg("CONTAINER", []string{}, "Container name or id")
	)

	cmd.Spec = "[-f | --force] CONTAINER..."

	cmd.Before = container.InitDatabase

	cmd.Action = func() {
		container.RemoveContainer(*containerID, *force)
	}
}
