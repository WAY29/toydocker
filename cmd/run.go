package cmd

import (
	"github.com/WAY29/toydocker/cgroups/subsystems"
	"github.com/WAY29/toydocker/container"
	"github.com/WAY29/toydocker/structs"
	cli "github.com/jawher/mow.cli"
)

func CmdRun(cmd *cli.Cmd) {
	var (
		tty         = cmd.BoolOpt("t tty", false, "Allocate a pseudo-TTY")
		interactive = cmd.BoolOpt("i interactive", false, "Keep STDIN open even if not attached")
		memory      = cmd.StringOpt("m memory", "1024m", "memory limit")
		cpushare    = cmd.StringOpt("cpushare", "1024", "cpushare limit")
		cpuset      = cmd.StringOpt("cpuset", "2", "cpuset limit")
		imagePath   = cmd.StringOpt("p path", "./images/busybox.tar", "Specifies the path of the image")
		volume      = cmd.StringsOpt("v volume", []string{}, "Bind mount a volume")
		detach      = cmd.BoolOpt("d detach", false, "Run container in background and print container ID")
	)

	var (
		command = cmd.StringsArg("COMMAND", []string{}, "command to run")
	)

	cmd.Spec = "[-t | --tty] [-i | --interactive] [-d | --detach] [-m=<memory limit> | --memory=<memory limit>] [--cpushare=<cpushare limit>] [--cpuset=<cpuset limit>] [-v | --volume]... (-p=<image path> | --path=<iamge path>) COMMAND..."

	cmd.Action = func() {
		resource := &subsystems.ResourceConfig{
			MemoryLimit: *memory,
			CpuShare:    *cpushare,
			CpuSet:      *cpuset,
		}

		cmdConfig := &structs.CmdConfig{
			Tty:         *tty,
			Interactive: *interactive,
			Detach:      *detach,
			ImagePath:   *imagePath,
			Volume:      *volume,
		}

		container.Run(cmdConfig, *command, resource)
	}
}
