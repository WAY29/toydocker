package subsystems

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type MemorySubSystem struct {
	cgroupPath string
}

// 返回subsystem名字
func (s *MemorySubSystem) Name() string {
	return "memory"
}

// 获取cgroup路径
func (s *MemorySubSystem) GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	if s.cgroupPath != "" {
		return s.cgroupPath, nil
	} else {
		cgroupPath, err := getCgroupPath(s.Name(), cgroupPath, true)
		s.cgroupPath = cgroupPath
		return cgroupPath, err
	}
}

// 设置某个cgroup在这个subsystem中的资源限制
func (s *MemorySubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	if subsysCgroupPath, err := s.GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "memory.oom_control"), []byte("1"), 0644); err != nil {
			return fmt.Errorf("set cgroup memory.oom_control error: %v", err)
		}
		if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "memory.limit_in_bytes"), []byte(res.MemoryLimit), 0644); err != nil {
			return fmt.Errorf("set cgroup memory.limit_in_bytes error: %v", err)
		}
	} else {
		return err
	}

	return nil
}

// 将进程添加到某个cgroup中
func (s *MemorySubSystem) Apply(cgroupPath string, pid int) error {
	if subsysCgroupPath, err := s.GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("set cgroup proc error: %v", err)
		}
	} else {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}
	return nil
}

// 移除某个cgroup
func (s *MemorySubSystem) Remove(cgroupPath string) error {
	if subsysCgroupPath, err := s.GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		return os.Remove(subsysCgroupPath)
	}
	return nil
}
