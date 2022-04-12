package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	cli "github.com/jawher/mow.cli"
	"github.com/sirupsen/logrus"
)

func getEnvsByPid(pid string) []string {
	path := fmt.Sprintf("/proc/%s/environ", pid)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Errorf("Read file %s error: %v", path, err)
		return nil
	}
	envs := strings.Split(string(content), "\u0000")
	return envs
}

func ExecContainer(container string, commandArray, envs []string) {
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
	containerEnvs := getEnvsByPid(pid)
	cmd.Env = append(os.Environ(), envs...)
	cmd.Env = append(cmd.Env, containerEnvs...)

	if err := cmd.Run(); err != nil {
		logrus.Error(err)
		cli.Exit(1)
	}

}
