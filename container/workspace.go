package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"syscall"

	"github.com/WAY29/toydocker/utils"
	cli "github.com/jawher/mow.cli"
	"github.com/sirupsen/logrus"
)

//Create a AUFS filesystem as container root workspace
func newWorkSpace(rootPath, ImagePath, containerID string, volumes []string) string {
	mntRootPath := path.Join(rootPath, "mnt")
	imageRootPath := path.Join(rootPath, "images")
	writeLayerRootPath := path.Join(rootPath, "write-layers")
	volumeRootPath := path.Join(rootPath, "volumes")

	if err := os.MkdirAll(rootPath, 0777); err != nil && !os.IsExist(err) {
		logrus.Error("Mkdir %s error: %v", rootPath, err)
		cli.Exit(1)
	}
	if err := os.MkdirAll(mntRootPath, 0777); err != nil && !os.IsExist(err) {
		logrus.Error("Mkdir %s error: %v", mntRootPath, err)
		cli.Exit(1)
	}
	if err := os.MkdirAll(imageRootPath, 0777); err != nil && !os.IsExist(err) {
		logrus.Error("Mkdir %s error: %v", imageRootPath, err)
		cli.Exit(1)
	}
	if err := os.MkdirAll(writeLayerRootPath, 0777); err != nil && !os.IsExist(err) {
		logrus.Error("Mkdir %s error: %v", writeLayerRootPath, err)
		cli.Exit(1)
	}
	if err := os.MkdirAll(volumeRootPath, 0777); err != nil && !os.IsExist(err) {
		logrus.Error("Mkdir %s error: %v", volumeRootPath, err)
		cli.Exit(1)
	}
	readonlyLayerPath := createReadOnlyLayer(imageRootPath, ImagePath)
	writeLayerPath := createWriteLayer(writeLayerRootPath, containerID)
	mntPath := createMountPoint(rootPath, mntRootPath, readonlyLayerPath, writeLayerPath, containerID)
	createVolumes(mntPath, volumes)

	return mntPath
}

func createReadOnlyLayer(imageRootPath, ImagePath string) string {
	// 判断镜像是否存在
	exist, err := utils.PathExists(ImagePath)
	if err != nil {
		logrus.Error(err)
		cli.Exit(1)
	} else if exist == false {
		logrus.Errorf("Image %s not exist", ImagePath)
		cli.Exit(1)
	}
	// 计算Image的hash,判断文件是否已经创建
	iamgeHash, err := utils.FileHash(ImagePath)
	if err != nil {
		logrus.Error(err)
		cli.Exit(1)
	}

	ImageDecompressionPath := path.Join(imageRootPath, iamgeHash)
	exist, err = utils.PathExists(ImageDecompressionPath)
	if err != nil {
		logrus.Error(err)
		cli.Exit(1)
	} else if exist == false {
		if err := os.Mkdir(ImageDecompressionPath, 0777); err != nil {
			logrus.Error("Mkdir %s error: %v", ImageDecompressionPath, err)
			cli.Exit(1)
		}
		if err = utils.Untar(ImagePath, ImageDecompressionPath); err != nil {
			logrus.Errorf("Untar error: %v", err)
			cli.Exit(1)
		}
	}

	// 处理镜像都解压到单独文件夹的情况
	files, err := ioutil.ReadDir(ImageDecompressionPath)

	if err != nil {
		logrus.Error(err)
		cli.Exit(1)
	}

	if len(files) == 1 {
		ImageDecompressionPath = path.Join(ImageDecompressionPath, files[0].Name())
	}

	return ImageDecompressionPath
}

func createWriteLayer(writeLayerRootPath, containerID string) string {
	writeLayerPath := path.Join(writeLayerRootPath, containerID)
	if err := os.Mkdir(writeLayerPath, 0777); err != nil {
		logrus.Error("Mkdir %s error: %v", writeLayerPath, err)
		cli.Exit(1)
	}
	return writeLayerPath
}

func createMountPoint(rootPath, mntRootPath, readonlyLayerPath, writeLayerPath, containerID string) string {
	mntPath := path.Join(mntRootPath, containerID)

	if err := os.MkdirAll(mntPath, 0777); err != nil {
		logrus.Error("Mkdir %s error: %v", mntPath, err)
		cli.Exit(1)
	}

	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	if err := syscall.Mount("none", mntPath, "aufs", uintptr(syscall.MS_NODEV), fmt.Sprintf("dirs=%s:%s", writeLayerPath, readonlyLayerPath)); err != nil {
		logrus.Error("Mount %s error: %v", mntPath, err)
		cli.Exit(1)
	}

	return mntPath
}

func createVolumes(mntPath string, volumes []string) {
	if len(volumes) == 0 {
		return
	}

	for _, volume := range volumes {
		volumeURLS, ok := volumeURLExtract(volume)
		if !ok {
			logrus.Warningf("Volume parameter[%s] input is invalid", volume)
		} else {
			mountVolume(mntPath, volumeURLS)
			logrus.Infof("Mount volume %s", volume)

		}
	}
}

func volumeURLExtract(volume string) ([]string, bool) {
	volumeURLS := strings.Split(volume, ":")
	if len(volumeURLS) == 2 && volumeURLS[0] != "" && volumeURLS[1] != "" {
		return volumeURLS, true
	}
	return nil, false
}

func mountVolume(mntPath string, volumeURLS []string) {
	parentPath := volumeURLS[0]
	if err := os.Mkdir(parentPath, 0777); err != nil && !os.IsExist(err) {
		logrus.Errorf("Mkdir parent dir %s error: %v", parentPath, err)
		cli.Exit(1)
	}
	containerVolumePath := path.Join(mntPath, volumeURLS[1])
	if err := os.Mkdir(containerVolumePath, 0777); err != nil && !os.IsExist(err) {
		logrus.Error("Mkdir container dir %s error: %v", containerVolumePath, err)
		cli.Exit(1)
	}

	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	if err := syscall.Mount("none", containerVolumePath, "aufs", uintptr(syscall.MS_NODEV), fmt.Sprintf("dirs=%s", parentPath)); err != nil {
		logrus.Error("Mount %s error: %v", containerVolumePath, err)
		cli.Exit(1)
	}
}

//Delete the AUFS filesystem while container exit
func deleteWorkSpace(rootPath, mntPath, containerID string, volumes []string) {
	deleteVolumes(mntPath, volumes)
	deleteMountPoint(rootPath, mntPath)
	deleteWriteLayer(rootPath, containerID)
}

func deleteMountPoint(rootPath string, mntPath string) {
	if err := syscall.Unmount(mntPath, 0); err != nil {
		logrus.Errorf("Unmount %s error: %v", mntPath, err)
		cli.Exit(1)
	}

	if err := os.RemoveAll(mntPath); err != nil {
		logrus.Error("Remove dir %s error: %v", mntPath, err)
		cli.Exit(1)
	}
}

func deleteWriteLayer(rootPath, containerID string) {
	writeLayerPath := path.Join(rootPath, "write-layers", containerID)

	if err := os.RemoveAll(writeLayerPath); err != nil {
		logrus.Errorf("Remove dir %s error %v", writeLayerPath, err)
	}
}

func deleteVolumes(mntPath string, volumes []string) {
	if len(volumes) == 0 {
		return
	}

	for _, volume := range volumes {
		volumeURLS, ok := volumeURLExtract(volume)
		if !ok {
			continue
		}
		unmountPath := path.Join(mntPath, volumeURLS[1])

		if err := syscall.Unmount(unmountPath, 0); err != nil {
			logrus.Errorf("Unmount %s error: %v", unmountPath, err)
		}
	}
}
