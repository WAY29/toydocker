package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"syscall"

	"github.com/WAY29/toydocker/utils"
	cli "github.com/jawher/mow.cli"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

//Create a AUFS filesystem as container root workspace
func newWorkSpace(rootPath, ImagePath, containerName string) string {
	mntRootPath := path.Join(rootPath, "mnt")
	imageRootPath := path.Join(rootPath, "images")
	writeLayerRootPath := path.Join(rootPath, "write-layers")

	if err := os.MkdirAll(rootPath, 0777); err != nil {
		logrus.Error("Mkdir %s error: %v", rootPath, err)
		cli.Exit(1)
	}
	if err := os.MkdirAll(mntRootPath, 0777); err != nil {
		logrus.Error("Mkdir %s error: %v", mntRootPath, err)
		cli.Exit(1)
	}
	if err := os.MkdirAll(imageRootPath, 0777); err != nil {
		logrus.Error("Mkdir %s error: %v", imageRootPath, err)
		cli.Exit(1)
	}
	if err := os.MkdirAll(writeLayerRootPath, 0777); err != nil {
		logrus.Error("Mkdir %s error: %v", writeLayerRootPath, err)
		cli.Exit(1)
	}

	readonlyLayerPath := createReadOnlyLayer(imageRootPath, ImagePath)
	writeLayerPath := createWriteLayer(writeLayerRootPath, containerName)
	mntPath := createMountPoint(rootPath, mntRootPath, readonlyLayerPath, writeLayerPath, containerName)

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

	logrus.Infof("ImageDecompressionPath: %s", ImageDecompressionPath)

	return ImageDecompressionPath
}

func createWriteLayer(writeLayerRootPath, containerName string) string {
	writeLayerPath := path.Join(writeLayerRootPath, containerName)
	if err := os.Mkdir(writeLayerPath, 0777); err != nil {
		logrus.Error("Mkdir %s error: %v", writeLayerPath, err)
		cli.Exit(1)
	}
	return writeLayerPath
}

func createMountPoint(rootPath, mntRootPath, readonlyLayerPath, writeLayerPath, containerName string) string {
	mntPath := path.Join(mntRootPath, containerName)

	if err := os.MkdirAll(mntPath, 0777); err != nil {
		logrus.Error("Mkdir %s error: %v", mntPath, err)
		cli.Exit(1)
	}

	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	defaultMountFlags := syscall.MS_NODEV
	if err := syscall.Mount("none", mntPath, "aufs", uintptr(defaultMountFlags), fmt.Sprintf("dirs=%s:%s", writeLayerPath, readonlyLayerPath)); err != nil {
		logrus.Error("Mkdir %s error: %v", mntPath, err)
		cli.Exit(1)
	}

	return mntPath
}

//Delete the AUFS filesystem while container exit
func deleteWorkSpace(rootPath, mntPath, containerName string) {
	deleteMountPoint(rootPath, mntPath)
	deleteWriteLayer(rootPath, containerName)
}

func deleteMountPoint(rootPath string, mntPath string) {
	if err := syscall.Unmount(mntPath, 0); err != nil {
		logrus.Error(err)
		cli.Exit(1)
	}

	if err := os.RemoveAll(mntPath); err != nil {
		logrus.Error("Remove dir %s error: %v", mntPath, err)
		cli.Exit(1)
	}
}

func deleteWriteLayer(rootPath, containerName string) {
	writeLayerPath := path.Join(rootPath, "write-layers", containerName)

	if err := os.RemoveAll(writeLayerPath); err != nil {
		log.Errorf("Remove dir %s error %v", writeLayerPath, err)
	}
}
