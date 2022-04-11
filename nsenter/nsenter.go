package nsenter

/*
#define _GNU_SOURCE
#include <errno.h>
#include <sched.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>
#include <unistd.h>

__attribute__((constructor)) void enter_namespace(void) {
	char *_TOYDOCKER_NS_PID;
	_TOYDOCKER_NS_PID = getenv("_TOYDOCKER_NS_PID");
	if (_TOYDOCKER_NS_PID) {
	} else {
		return;
	}
	char *_TOYDOCKER_NS_CMD;
	_TOYDOCKER_NS_CMD = getenv("_TOYDOCKER_NS_CMD");
	if (_TOYDOCKER_NS_CMD) {
	} else {
		return;
	}
	int i;
	char nspath[1024];
	char *namespaces[] = { "ipc", "uts", "net", "pid", "mnt" };

	for (i=0; i<5; i++) {
		sprintf(nspath, "/proc/%s/ns/%s", _TOYDOCKER_NS_PID, namespaces[i]);
		int fd = open(nspath, O_RDONLY);

		if (setns(fd, 0) == -1) {
		} else {
		}
		close(fd);
	}
	int res = system(_TOYDOCKER_NS_CMD);
	exit(0);
	return;
}
*/
import "C"
