package global

import (
	"sun-panel/lib/monitor"
)

const (
	SystemMonitor_CPU_INFO    = "CPU_INFO"
	SystemMonitor_MEMORY_INFO = "MEMORY_INFO"
	SystemMonitor_DISK_INFO   = "DISK_INFO"
)

type ModelSystemMonitor struct {
	CPUInfo           monitor.CPUInfo             `json:"cpuInfo"`
	DiskInfo          []monitor.DiskInfo          `json:"diskInfo"`
	NetIOCountersInfo []monitor.NetIOCountersInfo `json:"netIOCountersInfo"`
	MemoryInfo        monitor.MemoryInfo          `json:"memoryInfo"`
}
