package subsystems

var SubsystemsIns = map[string]SubSystem{
	"memory": &MemorySubSystem{cgroupPath: ""},
	"cpu":    &CpuSubSystem{cgroupPath: ""},
	"cpuset": &CpusetSubSystem{cgroupPath: ""},
}
