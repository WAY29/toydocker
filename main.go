package main

import (
	"os"

	"github.com/WAY29/toydocker/cmd"

	nested "github.com/antonfisher/nested-logrus-formatter"
	cli "github.com/jawher/mow.cli"
	"github.com/sirupsen/logrus"
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
		logrus.SetFormatter(&nested.Formatter{
			HideKeys:    true,
			FieldsOrder: []string{"component", "category"},
		})
		logrus.SetOutput(os.Stdout)
	}

	app.Command("init", "Init container process run user's process in container. Do not call it outside", cmd.CmdInit)
	app.Command("run", "Create a container with namespace and cgroups", cmd.CmdRun)
	app.Command("export", "Export a container's filesystem as a tar archive", cmd.CmdExport)
	app.Command("logs", "Fetch the logs of a container", cmd.CmdLogs)
	app.Command("ps", "List containers", cmd.CmdPS)

	app.Run(os.Args)
}
