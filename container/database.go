package container

import (
	"database/sql"
	"fmt"
	"path"
	"strings"

	"github.com/WAY29/toydocker/utils"
	cli "github.com/jawher/mow.cli"
	"github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"
)

type ContainerInfo struct {
	Id          int
	Pid         int
	ContainerId string
	Command     string
	CreateTime  string
	ImagePath   string
	Ports       string
	Status      string
	Name        string
}

const (
	PROC_STATUS_RUNNING  = "Running"
	PROC_STATUS_SLEEPING = "Sleeping"
	PROC_STATUS_STOP     = "Stop"
	PROC_STATUS_EXIT     = "Exited"
	PROC_STATUS_ZOMBIE   = "Zombie"
	PROC_STATUS_DEAD     = "Dead"
)

var (
	DB *sql.DB
)

func InitDatabase() {
	initDatabase(path.Join(ROOT_PATH, "database.db"))
}

func initDatabase(databasePath string) {
	var err error
	DB, err = sql.Open("sqlite", databasePath)
	if err != nil {
		logrus.Errorf("Open sqlite %s error: %v", databasePath, err)
		cli.Exit(1)
	}

	if exist, err := utils.PathExists(databasePath); err != nil {
		logrus.Error(err)
		cli.Exit(1)
	} else if !exist {
		createTable(databasePath)
	}
}

func createTable(databasePath string) {
	createTableQuery := strings.ReplaceAll(`CREATE TABLE 'ContainerInfomation' (
'id' INTEGER PRIMARY KEY AUTOINCREMENT,
'pid' INTEGER NULL,
'containerId' VARCHAR(128) NULL,
'command' VARCHAR(1024) NULL,
'createTime' VARCHAR(128) NULL,
'imagePath' VARCHAR(1024) NULL,
'ports' VARCHAR(2048) NULL,
'status' VARCHAR(32) NULL,
'name' VARCHAR(128) NULL);`, "'", "`")
	if _, err := runQuery(DB, createTableQuery); err != nil {
		logrus.Errorf("Create table ContainerInfomation error: %v", err)
		cli.Exit(1)
	}
}

func recordContainerInfo(containerInfo *ContainerInfo) {
	insertQuery := `INSERT INTO ContainerInfomation(pid, containerId, command, createTime, imagePath, ports, status, name) values(?,?,?,?,?,?,?,?)`
	_, err := runQuery(DB, insertQuery, containerInfo.Pid, containerInfo.ContainerId, containerInfo.Command, containerInfo.CreateTime, containerInfo.ImagePath, containerInfo.Ports, containerInfo.Status, containerInfo.Name)
	if err != nil {
		logrus.Errorf("Record ContainerInfo error: %v", err)
		cli.Exit(1)
	}
}

func ListContainerInfo() []ContainerInfo {
	var containerInfo *ContainerInfo
	containerInfos := make([]ContainerInfo, 0)

	selectQuery := `SELECT pid, containerId, command, createTime, imagePath, ports, status, name FROM ContainerInfomation ORDER BY containerId`
	rows, err := DB.Query(selectQuery)
	if err != nil {
		logrus.Errorf("List ContainerInfo error: %v", err)
		cli.Exit(1)
	}

	for rows.Next() {
		containerInfo = &ContainerInfo{}
		err = rows.Scan(&containerInfo.Pid, &containerInfo.ContainerId, &containerInfo.Command, &containerInfo.CreateTime, &containerInfo.ImagePath, &containerInfo.Ports, &containerInfo.Status, &containerInfo.Name)
		if err != nil {
			logrus.Errorf("Get ContainerInfo from db error: %v", err)
			cli.Exit(1)
		}

		containerInfos = append(containerInfos, *containerInfo)
	}

	return containerInfos
}

func updateContainerProperty(containerID, property, value string) {
	updateQuery := fmt.Sprintf(`UPDATE ContainerInfomation SET %s=?`, property)
	_, err := runQuery(DB, updateQuery, value)
	if err != nil {
		logrus.Errorf("Update Container[%s] status[%s] error: %v", containerID, value, err)
		cli.Exit(1)
	}
}

func updateContainerStatus(containerID, status string) {
	updateContainerProperty(containerID, "status", status)
}

func updateContainerPID(containerID, pid string) {
	updateContainerProperty(containerID, "pid", pid)
}

func findContainerProperty(container string, properties ...string) []string {
	propertieStr := strings.Join(properties, ",")

	selectQuery := fmt.Sprintf("SELECT %s FROM ContainerInfomation WHERE containerId LIKE (? ||'%%') or name=?", propertieStr)
	rows, err := DB.Query(selectQuery, container, container)
	if err != nil {
		logrus.Errorf("List ContainerInfo error: %v", err)
		cli.Exit(1)
	}

	values := make([]string, len(properties))
	scans := make([]interface{}, len(properties))
	for i := range scans {
		scans[i] = &values[i]
	}
	theOnlyFlag := false

	for rows.Next() {
		if !theOnlyFlag {
			theOnlyFlag = true
		} else {
			logrus.Errorf("Ambiguous container id or name: %s", container)
			cli.Exit(1)
		}
		err = rows.Scan(scans...)
		if err != nil {
			logrus.Errorf("Get ContainerInfo from db error: %v", err)
			cli.Exit(1)
		}
	}

	return values

}

func findContainerPID(container string) string {
	if result := findContainerProperty(container, "pid"); len(result) > 0 {
		return result[0]
	}
	return ""
}

func findContainerID(container string) string {
	if result := findContainerProperty(container, "containerID"); len(result) > 0 {
		return result[0]
	}
	return ""
}

func findContainerIDAndPID(container string) (string, string) {
	if result := findContainerProperty(container, "containerID", "pid"); len(result) > 0 {
		return result[0], result[1]
	}
	return "", ""
}

func runQuery(db *sql.DB, query string, args ...interface{}) (sql.Result, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}
