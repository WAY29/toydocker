package container

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/WAY29/toydocker/cgroups"
	"github.com/WAY29/toydocker/cgroups/subsystems"
	"github.com/WAY29/toydocker/structs"
	"github.com/WAY29/toydocker/utils"
	"github.com/sirupsen/logrus"
)

func newPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}

// fork自身并创建namespace隔离
func newParentProcess(ttyFlag, interactiveFlag, detachFlag bool, commandArray []string) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := newPipe()
	if err != nil {
		logrus.Errorf("New pipe error %v", err)
		return nil, nil
	}

	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET,
	}
	if !detachFlag {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if interactiveFlag {
		cmd.Stdin = os.Stdin
	}

	// cmd.Dir = BUSYBOX_IMAGE_DIR
	cmd.ExtraFiles = []*os.File{readPipe}

	return cmd, writePipe
}

func Run(cmdConfig *structs.CmdConfig, commandArray []string, resource *subsystems.ResourceConfig) {
	parent, writePipe := newParentProcess(cmdConfig.Tty, cmdConfig.Interactive, cmdConfig.Detach, commandArray)

	containerID := utils.RandStr(12)
	logrus.Infof("containterID: %s", containerID)

	// 创建workspace
	mntPath := newWorkSpace(ROOT_PATH, cmdConfig.ImagePath, containerID, cmdConfig.Volume)
	// if !cmdConfig.Detach {
	// 	defer deleteWorkSpace(ROOT_PATH, mntPath, containerID, cmdConfig.Volume)
	// }
	// 设置新的文件系统根目录
	parent.Dir = mntPath

	// 创建cgroup manager，并调用set和apply设置资源限制
	cgroupManager := cgroups.NewCgroupManager(containerID)
	// if !cmdConfig.Detach {
	// 	defer cgroupManager.Destroy()
	// }
	// 设置cgroup资源限制
	cgroupManager.Set(resource)

	// 启动容器进程
	if err := parent.Start(); err != nil {
		logrus.Error(err)
	}
	pid := parent.Process.Pid
	logrus.Infof("Pid: %d", pid)

	// 添加pid到cgroup,使用资源限制
	cgroupManager.Apply(pid)

	// 设置命令参数
	command := sendInitCommand(commandArray, writePipe)

	// 记录容器数据
	if cmdConfig.Name == "" {
		cmdConfig.Name = containerID
	}
	recordContainerInfo(&ContainerInfo{
		Pid:         parent.Process.Pid,
		ContainerId: containerID,
		Command:     command,
		CreateTime:  time.Now().Format("2006-01-02 15:04:05"),
		ImagePath:   cmdConfig.ImagePath,
		Ports:       "",
		Status:      RUNNING,
		Name:        cmdConfig.Name,
	})

	// 若非放到后台运行则等待进程
	if !cmdConfig.Detach {
		parent.Wait()
	} else {
		// 否则输出containerID
		fmt.Println(containerID)
	}

}

// 使用管道传输命令参数
func sendInitCommand(comArray []string, writePipe *os.File) string {
	command := strings.Join(comArray, " ")
	logrus.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()

	return command
}
