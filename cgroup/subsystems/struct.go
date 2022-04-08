package subsystems

// 用于传递资源限制配置的结构体
type ResourceConfig struct {
	// 内存
	MemoryLimit string
	// cpu时间片权重
	CpuShare string
	// cpu核心数
	CpuSet string
}

type SubSystem interface {
	// 返回subsystem名字
	Name() string
	// 设置某个cgroup在这个subsystem中的资源限制
	Set(path string, res *ResourceConfig) error
	// 获取cgroup路径
	GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error)
	// 将进程添加到某个cgroup中
	Apply(path string, pid int) error
	// 移除某个cgroup
	Remove(path string) error
}
