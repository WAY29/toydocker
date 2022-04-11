package container

import (
	"os"
	"os/exec"
	"strings"

	cli "github.com/jawher/mow.cli"
	"github.com/sirupsen/logrus"
)

func ExecContainer(container string, commandArray []string) {
	pid := findContainerPID(container)
	if pid == "" {
		logrus.Warningf("No such container: %s", container)
		cli.Exit(1)
	}
	commandStr := strings.Join(commandArray, " ")
	cmd := exec.Command("/proc/self/exe", "exec", "nothing", "nothing")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	os.Setenv(ENV_EXEC_PID, pid)
	os.Setenv(ENV_EXEC_CMD, commandStr)

	if err := cmd.Run(); err != nil {
		logrus.Error(err)
		cli.Exit(1)
	}

}
