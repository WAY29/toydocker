package container

import (
	"os"
	"syscall"

	"github.com/google/shlex"
	log "github.com/sirupsen/logrus"
)

func RunContainerinitProcess(command string) {
	var (
		err  error
		argv []string
	)

	log.Infof("command %s", command)

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	// argv := []string{command}
	if argv, err = shlex.Split(command); err != nil {
		log.Error(err)
		os.Exit(1)
	}

	if err = syscall.Exec(argv[0], argv, os.Environ()); err != nil {
		log.Error(err)
	}

}
