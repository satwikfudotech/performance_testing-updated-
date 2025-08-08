package utils

import (
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func GetCPUUsage() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	if len(percent) > 0 {
		return percent[0]
	}
	return 0
}

func GetRAMUsage() float64 {
	v, _ := mem.VirtualMemory()
	return float64(v.Used) / (1024 * 1024) // MB
}

func GetGoroutineCount() int {
	return runtime.NumGoroutine()
}
