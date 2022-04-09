package subsystems

var SubsystemsIns = []SubSystem{
	&MemorySubSystem{cgroupPath: ""},
	&CpuSubSystem{cgroupPath: ""},
	&CpusetSubSystem{cgroupPath: ""},
}
