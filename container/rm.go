package container

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/WAY29/toydocker/cgroups"
	cli "github.com/jawher/mow.cli"
	"github.com/sirupsen/logrus"
)

func RemoveContainer(containers []string, force bool) {
	for _, container := range containers {
		results := findContainerProperties(container, "containerId", "pid", "status", "volumes")

		containerID, pid, status, volumeStr := results[0], results[1], results[2], results[3]
		if pid == "" || containerID == "" {
			logrus.Warningf("No such container: %s", container)
			cli.Exit(1)
		}
		volumes := strings.Split(volumeStr, "||")

		if status != PROC_STATUS_EXIT {
			if !force {
				logrus.Warningf("Can't remove %s container", status)
				cli.Exit(1)
			} else {
				StopContainer([]string{containerID})
			}
		}

		// 删除cgroup
		cgroupManager := cgroups.NewCgroupManager(containerID)
		cgroupManager.Destroy()

		// 删除workspace
		mntPath := path.Join(ROOT_PATH, "mnt", containerID)
		deleteWorkSpace(ROOT_PATH, mntPath, containerID, volumes)

		// 删除日志
		logPath := path.Join(ROOT_PATH, "logs", containerID)
		if err := os.RemoveAll(logPath); err != nil && !os.IsNotExist(err) {
			logrus.Errorf("Remove log error: %v", err)
			cli.Exit(1)
		}

		// 删除记录
		if err := deleteContainerInfo(containerID); err != nil {
			logrus.Errorf("Delete container info from db error: %v", err)
			cli.Exit(1)
		}

		fmt.Printf("%s\n", containerID)
	}
}
