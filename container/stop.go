package container

import (
	"fmt"
	"strconv"
	"syscall"

	cli "github.com/jawher/mow.cli"
	"github.com/sirupsen/logrus"
)

func StopContainer(containers []string) {
	for _, container := range containers {
		containerID, pid := findContainerIDAndPID(container)
		if pid == "" || containerID == "" {
			logrus.Warningf("No such container: %s", container)
			cli.Exit(1)
		}

		pidInt, err := strconv.Atoi(pid)
		if err != nil {
			logrus.Errorf("Conver pid from string to int error: %v", err)
			cli.Exit(1)
		}
		if err = syscall.Kill(pidInt, syscall.SIGTERM); err != nil {
			logrus.Errorf("Stop pid %d: %v", pidInt, err)
			cli.Exit(1)
		}

		updateContainerPID(containerID, "-1")

		fmt.Printf("%s\n", pid)
	}
}
