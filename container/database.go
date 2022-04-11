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
	RUNNING = "Running"
	STOP    = "Stop"
	EXIT    = "Exit"
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

func updateContainerStatus(containerID, status string) {
	updateQuery := `UPDATE ContainerInfomation SET status=?`
	_, err := runQuery(DB, updateQuery, status)
	if err != nil {
		logrus.Errorf("Update Container[%s] status[%s] error: %v", containerID, status, err)
		cli.Exit(1)
	}
}

func findContainerInfo(container, property string) string {
	selectQuery := fmt.Sprintf("SELECT %s FROM ContainerInfomation WHERE containerId LIKE (? ||'%%') or name=?", property)
	rows, err := DB.Query(selectQuery, container, container)
	if err != nil {
		logrus.Errorf("List ContainerInfo error: %v", err)
		cli.Exit(1)
	}

	var value string = ""
	theOnlyFlag := false

	for rows.Next() {
		if !theOnlyFlag {
			theOnlyFlag = true
		} else {
			logrus.Errorf("Ambiguous container id or name: %s", container)
			cli.Exit(1)
		}
		err = rows.Scan(&value)
		if err != nil {
			logrus.Errorf("Get ContainerInfo from db error: %v", err)
			cli.Exit(1)
		}
	}

	return value

}

func findContainerPID(container string) string {
	return findContainerInfo(container, "pid")
}

func findContainerID(container string) string {
	return findContainerInfo(container, "containerID")

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
