package container

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/creack/pty"
	"golang.org/x/term"

	"github.com/WAY29/toydocker/cgroups"
	"github.com/WAY29/toydocker/cgroups/subsystems"
	"github.com/WAY29/toydocker/structs"
	"github.com/WAY29/toydocker/utils"
	cli "github.com/jawher/mow.cli"
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
func newParentProcess(ttyFlag, interactiveFlag, detachFlag bool, commandArray []string, containerID string) (*exec.Cmd, *os.File, *io.Writer) {
	readPipe, writePipe, err := newPipe()
	if err != nil {
		logrus.Errorf("New pipe error %v", err)
		return nil, nil, nil
	}

	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET,
	}

	// 记录标准输出错误到日志文件中
	var (
		stdoutWriter io.Writer
		stderrWriter io.Writer
	)
	logPath := path.Join(ROOT_PATH, "logs", containerID)
	logFilePath := path.Join(logPath, "container.log")
	checkErr(utils.MkdirAll(logPath, 0622), "Mkdir %s error", logPath)
	logFile, err := os.Create(logFilePath)
	if err != nil {
		logrus.Errorf("Create log %s error: %v", logFilePath, err)
		cli.Exit(1)
	}
	if !detachFlag {
		stdoutWriter = io.MultiWriter(os.Stdout, logFile)
		stderrWriter = io.MultiWriter(os.Stderr, logFile)
	} else {
		stdoutWriter = logFile
		stderrWriter = logFile
	}

	// 重定向输入输出错误
	if !ttyFlag {
		cmd.Stdout = stdoutWriter
		cmd.Stderr = stderrWriter
		if interactiveFlag {
			cmd.Stdin = os.Stdin
		}
	}

	cmd.ExtraFiles = []*os.File{readPipe}

	return cmd, writePipe, &stdoutWriter
}

func Run(cmdConfig *structs.CmdConfig, commandArray []string, resource *subsystems.ResourceConfig) {
	// 不允许相同名字的容器
	if !findContainerNameNotExist(cmdConfig.Name) {
		logrus.Errorf("A container with the same name[%s] exists", cmdConfig.Name)
		cli.Exit(1)
	}
	// 创建日志目录
	logRootPath := path.Join(ROOT_PATH, "logs")
	checkErr(utils.MkdirAll(logRootPath, 0622), "Mkdir %s error", logRootPath)

	// 生成containerID
	containerID := utils.RandStr(12)
	logrus.Infof("containterID: %s", containerID)

	parent, writePipe, stdoutWriter := newParentProcess(cmdConfig.Tty, cmdConfig.Interactive, cmdConfig.Detach, commandArray, containerID)

	// 创建workspace
	mntPath := newWorkSpace(ROOT_PATH, cmdConfig.ImagePath, containerID, cmdConfig.Volume)

	// 设置新的文件系统根目录
	parent.Dir = mntPath

	// 创建cgroup manager，并调用set和apply设置资源限制
	cgroupManager := cgroups.NewCgroupManager(containerID)
	if resource.MemoryLimit == "" {
		delete(subsystems.SubsystemsIns, "memory")
	}
	if resource.CpuShare == "" {
		delete(subsystems.SubsystemsIns, "cpu")
	}
	if resource.CpuSet == "" {
		delete(subsystems.SubsystemsIns, "cpuset")
	}

	// 设置cgroup资源限制
	cgroupManager.Set(resource)

	// 启动容器进程
	if cmdConfig.Tty {
		ptmx, err := pty.Start(parent)
		if err != nil {
			logrus.Error(err)
			cli.Exit(1)
		}
		defer func() { _ = ptmx.Close() }() // Best effort.
		// Make sure to close the pty at the end.
		// Handle pty size.
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGWINCH)
		go func() {
			for range ch {
				if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
					log.Printf("error resizing pty: %s", err)
				}
			}
		}()
		ch <- syscall.SIGWINCH                        // Initial resize.
		defer func() { signal.Stop(ch); close(ch) }() // Cleanup signals when done.

		// Set stdin in raw mode.
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			logrus.Error(err)
			cli.Exit(1)
		}
		defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.
		// Copy stdin to the pty and the pty to stdout.
		go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
		go func() { _, _ = io.Copy(*stdoutWriter, ptmx) }()
	} else if err := parent.Start(); err != nil {
		logrus.Error(err)
		cli.Exit(1)
	}

	pid := parent.Process.Pid

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
		Volumes:     cmdConfig.Volume,
		Ports:       "",
		Status:      PROC_STATUS_RUNNING,
		Name:        cmdConfig.Name,
	})

	// 若非放到后台运行则等待进程
	if !cmdConfig.Detach {
		defer updateContainerStatus(containerID, PROC_STATUS_EXIT)
		parent.Wait()
	} else {
		// 否则输出containerID
		fmt.Println(containerID)
	}

}

// 使用管道传输命令参数
func sendInitCommand(commandArray []string, writePipe *os.File) string {
	command := strings.Join(commandArray, " ")
	writePipe.WriteString(command)
	writePipe.Close()

	return command
}
