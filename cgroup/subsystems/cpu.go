package subsystems

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type CpuSubSystem struct {
	cgroupPath string
}

// 获取cgroup路径
func (s *CpuSubSystem) GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	if s.cgroupPath != "" {
		return s.cgroupPath, nil
	} else {
		cgroupPath, err := getCgroupPath(s.Name(), cgroupPath, true)
		s.cgroupPath = cgroupPath
		return cgroupPath, err
	}
}

// 设置某个cgroup在这个subsystem中的资源限制
func (s *CpuSubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	if subsysCgroupPath, err := s.GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		if res.CpuShare != "" {
			if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "cpu.shares"), []byte(res.CpuShare), 0644); err != nil {
				return fmt.Errorf("set cgroup cpu share fail %v", err)
			}
		}
		return nil
	} else {
		return err
	}
}

// 移除某个cgroup
func (s *CpuSubSystem) Remove(cgroupPath string) error {
	if subsysCgroupPath, err := s.GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		return os.RemoveAll(subsysCgroupPath)
	} else {
		return err
	}
}

// 将进程添加到某个cgroup中
func (s *CpuSubSystem) Apply(cgroupPath string, pid int) error {
	if subsysCgroupPath, err := s.GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("set cgroup proc fail %v", err)
		}
		return nil
	} else {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}
}

// 返回subsystem名字
func (s *CpuSubSystem) Name() string {
	return "cpu"
}
