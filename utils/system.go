package utils

import (
	"os"

	"github.com/shirou/gopsutil/process"
)

type SystemInfomation struct {
	Mem uint64
	CPU float64
}

func GetSystemInformation() *SystemInfomation {
	// get memory
	getMem := func(pro *process.Process) uint64 {
		processMemory, err := pro.MemoryInfo()
		if err != nil {
			return 0
		}
		return processMemory.RSS
	}
	curProcess, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return nil
	}
	mem := getMem(curProcess)
	cpu, err := curProcess.CPUPercent()
	if err == nil {
		return &SystemInfomation{Mem: mem, CPU: 0}
	}
	return &SystemInfomation{Mem: mem, CPU: cpu}
}
