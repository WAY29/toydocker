package container

import (
	"fmt"
	"path"

	"github.com/WAY29/toydocker/utils"
	"github.com/sirupsen/logrus"
)

func ExportContainer(outputPath, container string) {
	containerID := findContainerID(container)

	mntURL := path.Join(ROOT_PATH, "mnt", containerID)
	if err := utils.Tar(mntURL, outputPath); err != nil {
		logrus.Errorf("Tar directory %s error: %v", mntURL, err)
	}

	fmt.Printf("%s\n", containerID)
}
