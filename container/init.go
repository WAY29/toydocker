package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	cli "github.com/jawher/mow.cli"
	log "github.com/sirupsen/logrus"
)

func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Errorf("init read pipe error %v", err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}

func RunContainerInitProcess() error {
	var (
		err error
	)

	argv := readUserCommand()
	if argv == nil || len(argv) == 0 {
		return fmt.Errorf("Run container get user command error, cmdArray is nil")
	}

	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	// 临时的寻找可执行程序方案，后续存在镜像后需要修改
	if binary, err := exec.LookPath(argv[0]); err == nil {
		argv[0] = binary
	} else {
		log.Error(err)
		cli.Exit(1)
	}

	if err = syscall.Exec(argv[0], argv, os.Environ()); err != nil {
		log.Error(err)
	}
	return nil

}
