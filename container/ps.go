package container

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"text/tabwriter"

	cli "github.com/jawher/mow.cli"
	"github.com/sirupsen/logrus"
)

var (
	STATUS_MAP = map[string]string{
		"R": PROC_STATUS_RUNNING,
		"S": PROC_STATUS_SLEEPING,
		"D": PROC_STATUS_SLEEPING,
		"T": PROC_STATUS_STOP,
		"Z": PROC_STATUS_ZOMBIE,
		"X": PROC_STATUS_DEAD,
	}
)

func PS() {
	ListContainers()
}

func getProcStatus(pid int) string {
	content, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		if os.IsNotExist(err) {
			return PROC_STATUS_EXIT
		} else {
			logrus.Errorf("Get process %s status error: %v", pid, err)
			cli.Exit(1)
		}
	}
	contentSplited := bytes.Split(content, []byte{0x20})
	status := contentSplited[2]
	if status, ok := STATUS_MAP[string(status)]; ok {
		return status
	}

	return "Unknown"

}

func ListContainers() {
	var status string

	containerInfos := listContainerInfo()

	w := tabwriter.NewWriter(os.Stdout, 5, 30, 1, ' ', 0)
	fmt.Fprintln(w, "CONTAINER_ID\tPID\tIMAGE\tCOMMAND\tCREATED_TIME\tSTATUS\tPORTS\tNAMES")

	for _, containerInfo := range containerInfos {
		status = getProcStatus(containerInfo.Pid)
		pid := strconv.Itoa(containerInfo.Pid)
		if pid == "-1" {
			pid = ""
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", containerInfo.ContainerId, pid, containerInfo.ImagePath, containerInfo.Command, containerInfo.CreateTime, status, containerInfo.Ports, containerInfo.Name)

		if status != containerInfo.Status {
			updateContainerStatus(containerInfo.ContainerId, status)
		}
	}

	w.Flush()
}
