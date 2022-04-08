package container

import (
	"os"
	"os/exec"
	"strings"
	"syscall"

	cgroups "github.com/WAY29/toydocker/cgroup"
	"github.com/WAY29/toydocker/cgroup/subsystems"
	_ "github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func newPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}

// fork自身并创建namespace隔离
func newParentProcess(ttyFlag, interactiveFlag bool, commandArray []string) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := newPipe()
	if err != nil {
		log.Errorf("New pipe error %v", err)
		return nil, nil
	}

	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET,
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if interactiveFlag {
		cmd.Stdin = os.Stdin
	}

	cmd.Dir = BUSYBOX_IMAGE_DIR
	cmd.ExtraFiles = []*os.File{readPipe}

	return cmd, writePipe
}

func Run(tty, interactive bool, commandArray []string, resource *subsystems.ResourceConfig) {
	parent, writePipe := newParentProcess(tty, interactive, commandArray)
	// 创建cgroup manager，并调用set和apply设置资源限制
	// uuid := uuid.New().String()
	uuid := "my-cgroup"
	log.Infof("cgroup name %s", uuid)

	cgroupManager := cgroups.NewCgroupManager(uuid)
	defer cgroupManager.Destroy()
	// 设置cgroup资源限制
	cgroupManager.Set(resource)

	// 启动容器进程
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	log.Infof("pid %d", parent.Process.Pid)

	// 添加pid到cgroup
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(commandArray, writePipe)
	parent.Wait()
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
