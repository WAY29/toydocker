package main

import (
	"os"

	"github.com/WAY29/toydocker/cmd"
	"github.com/sirupsen/logrus"

	cli "github.com/jawher/mow.cli"
)

const (
	__version__ = "1.0.0"
)

var (
	app *cli.Cli
)

func main() {
	app = cli.App("toydocker", "toydocker is a simple container runtime implementation.")

	app.Before = func() {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(os.Stdout)
	}

	app.Command("init", "Init container process run user's process in container. Do not call it outside", cmd.CmdInit)
	app.Command("run", "Create a container with namespace and cgroups", cmd.CmdRun)
	app.Command("export", "Export a container's filesystem as a tar archive", cmd.CmdExport)

	app.Run(os.Args)
}
