package utils

import (
	"os"
	"os/exec"
	"runtime"

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

// 打开网址
func OpenUrl(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
