package cmd

import (
	"github.com/WAY29/toydocker/cgroup/subsystems"
	"github.com/WAY29/toydocker/container"
	cli "github.com/jawher/mow.cli"
)

func CmdRun(cmd *cli.Cmd) {
	var (
		tty         = cmd.BoolOpt("t tty", false, "Allocate a pseudo-TTY")
		interactive = cmd.BoolOpt("i interactive", false, "Keep STDIN open even if not attached")
	)

	var (
		command  = cmd.StringsArg("COMMAND", []string{}, "command to run")
		memory   = cmd.StringOpt("m memory", "1024m", "memory limit")
		cpushare = cmd.StringOpt("cpushare", "1024", "cpushare limit")
		cpuset   = cmd.StringOpt("cpuset", "2", "cpuset limit")
	)

	cmd.Spec = "[-t | --tty] [-i | --interactive] [-m=<memory limit> | --memory=<memory limit>] [--cpushare=<cpushare limit>] [--cpuset=<cpuset limit>] COMMAND..."

	cmd.Action = func() {
		resource := &subsystems.ResourceConfig{
			MemoryLimit: *memory,
			CpuShare:    *cpushare,
			CpuSet:      *cpuset,
		}

		container.Run(*tty, *interactive, *command, resource)
	}
}
