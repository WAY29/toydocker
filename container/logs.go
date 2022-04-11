package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	cli "github.com/jawher/mow.cli"
	"github.com/sirupsen/logrus"
)

func ShowContainerLogs(container string) {
	containerID := findContainerID(container)
	if containerID == "" {
		logrus.Warningf("No such container: %s", container)
		cli.Exit(1)
	}

	logFilePath := path.Join(ROOT_PATH, "logs", containerID, "container.log")
	if content, err := ioutil.ReadFile(logFilePath); err != nil {
		logrus.Errorf("Read log %s error: %v", logFilePath, err)
		cli.Exit(1)
	} else {
		fmt.Fprint(os.Stdout, string(content))
	}

}
