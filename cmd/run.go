package cmd

import (
	"github.com/WAY29/toydocker/cgroups/subsystems"
	"github.com/WAY29/toydocker/container"
	"github.com/WAY29/toydocker/structs"
	cli "github.com/jawher/mow.cli"
)

func CmdRun(cmd *cli.Cmd) {
	var (
		name        = cmd.StringOpt("name", "", "Container name")
		tty         = cmd.BoolOpt("t tty", false, "Allocate a pseudo-TTY")
		interactive = cmd.BoolOpt("i interactive", false, "Keep STDIN open even if not attached")
		detach      = cmd.BoolOpt("d detach", false, "Run container in background and print container ID")
		memory      = cmd.StringOpt("m memory", "", "memory limit")
		cpushare    = cmd.StringOpt("cpushare", "", "cpushare limit")
		cpuset      = cmd.StringOpt("cpuset", "", "cpuset limit")
		imagePath   = cmd.StringOpt("p path", "./images/busybox.tar", "Specifies the path of the image")
		volume      = cmd.StringsOpt("v volume", []string{}, "Bind mount a volume")
		envs        = cmd.StringsOpt("e env", []string{}, "Set environment variables")
	)

	var (
		command = cmd.StringsArg("COMMAND", []string{}, "Command to run")
	)

	cmd.Spec = "[--name] [-t | --tty] [-i | --interactive] [-d | --detach] [-m=<memory limit> | --memory=<memory limit>] [--cpushare=<cpushare limit>] [--cpuset=<cpuset limit>] [-v | --volume]... [-e | --env]... (-p=<image path> | --path=<iamge path>) COMMAND..."

	cmd.Before = container.InitDatabase

	cmd.Action = func() {
		resource := &subsystems.ResourceConfig{
			MemoryLimit: *memory,
			CpuShare:    *cpushare,
			CpuSet:      *cpuset,
		}

		cmdConfig := &structs.CmdRunConfig{
			Name:        *name,
			Tty:         *tty,
			Interactive: *interactive,
			Detach:      *detach,
			ImagePath:   *imagePath,
			Volume:      *volume,
			Env:         *envs,
		}

		container.Run(cmdConfig, *command, resource)
	}
}
