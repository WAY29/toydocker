package subsystems

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type CpusetSubSystem struct {
	cgroupPath string
}

// 获取cgroup路径
func (s *CpusetSubSystem) GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	if s.cgroupPath != "" {
		return s.cgroupPath, nil
	} else {
		cgroupPath, err := getCgroupPath(subsystem, cgroupPath, true)
		s.cgroupPath = cgroupPath
		return cgroupPath, err
	}
}

// 设置某个cgroup在这个subsystem中的资源限制
func (s *CpusetSubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	var (
		content          []byte
		err              error
		subsysCgroupPath string
	)
	cgroupMountPointPath := FindCgroupMountpoint(s.Name())

	if subsysCgroupPath, err = s.GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		defaultMemsPath := path.Join(cgroupMountPointPath, "cpuset.mems")
		content, err = ioutil.ReadFile(defaultMemsPath)
		if err != nil {
			return fmt.Errorf("Set cgroup cpuset.mems error: %v", err)
		}
		if err = ioutil.WriteFile(path.Join(subsysCgroupPath, "cpuset.mems"), []byte(content), 0644); err != nil {
			return fmt.Errorf("Set cgroup cpuset.mems error: %v", err)
		}

		if res.CpuSet != "" {
			if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "cpuset.cpus"), []byte(res.CpuSet), 0644); err != nil {
				return fmt.Errorf("Set cgroup cpuset.cpus error: %v", err)
			}
		} else {
			defaultCpusPath := path.Join(cgroupMountPointPath, "cpuset.cpus")
			content, err = ioutil.ReadFile(defaultCpusPath)
			if err != nil {
				return fmt.Errorf("Set cgroup cpuset.cpus error: %v", err)
			}
			if err = ioutil.WriteFile(path.Join(subsysCgroupPath, "cpuset.cpus"), []byte(content), 0644); err != nil {
				return fmt.Errorf("Set cgroup cpuset.cpus error: %v", err)
			}
		}
		return nil
	} else {
		return err
	}
}

// 移除某个cgroup
func (s *CpusetSubSystem) Remove(cgroupPath string) error {
	if subsysCgroupPath, err := s.GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		return os.RemoveAll(subsysCgroupPath)
	} else {
		return err
	}
}

// 将进程添加到某个cgroup中
func (s *CpusetSubSystem) Apply(cgroupPath string, pid int) error {
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
func (s *CpusetSubSystem) Name() string {
	return "cpuset"
}
