package container

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func PS() {
	ListContainers()
}

func ListContainers() {
	containerInfos := ListContainerInfo()

	w := tabwriter.NewWriter(os.Stdout, 5, 30, 1, ' ', 0)
	fmt.Fprintln(w, "CONTAINER_ID\tPID\tIMAGE\tCOMMAND\tCREATED_TIME\tSTATUS\tPORTS\tNAMES")

	for _, containerInfo := range containerInfos {
		fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s\n", containerInfo.ContainerId, containerInfo.Pid, containerInfo.ImagePath, containerInfo.Command, containerInfo.CreateTime, containerInfo.Status, containerInfo.Ports, containerInfo.Name)
	}

	w.Flush()
}
