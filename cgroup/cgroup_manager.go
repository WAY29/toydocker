package cgroups

import (
	"github.com/WAY29/toydocker/cgroup/subsystems"
	log "github.com/sirupsen/logrus"
)

type CgroupManager struct {
	// cgroup在hierarchy中的路径 相当于创建的cgroup目录相对于root cgroup目录的路径
	Path string
	// 资源配置
	Resource *subsystems.ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}

// 将进程pid加入到这个cgroup中
func (c *CgroupManager) Apply(pid int) error {
	log.Infof("start apply")
	for _, subSysIns := range subsystems.SubsystemsIns {
		if err := subSysIns.Apply(c.Path, pid); err != nil {
			log.Error(err)
		}
	}
	return nil
}

// 设置cgroup资源限制
func (c *CgroupManager) Set(res *subsystems.ResourceConfig) error {
	log.Infof("start set")
	for _, subSysIns := range subsystems.SubsystemsIns {
		if err := subSysIns.Set(c.Path, res); err != nil {
			log.Error(err)
		}
	}
	return nil
}

//释放cgroup
func (c *CgroupManager) Destroy() error {
	log.Infoln("cgroup destroy")
	for _, subSysIns := range subsystems.SubsystemsIns {
		if err := subSysIns.Remove(c.Path); err != nil {
			log.Warnf("remove cgroup fail %v", err)
		}
	}
	return nil
}
