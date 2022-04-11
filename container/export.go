package container

import (
	"fmt"
	"path"

	"github.com/WAY29/toydocker/utils"
	cli "github.com/jawher/mow.cli"
	"github.com/sirupsen/logrus"
)

func ExportContainer(outputPath, container string) {
	containerID := findContainerID(container)
	if containerID == "" {
		logrus.Warningf("No such container: %s", container)
		cli.Exit(1)
	}

	mntURL := path.Join(ROOT_PATH, "mnt", containerID)
	if err := utils.Tar(mntURL, outputPath); err != nil {
		logrus.Errorf("Tar directory %s error: %v", mntURL, err)
	}

	fmt.Printf("%s\n", containerID)
}
