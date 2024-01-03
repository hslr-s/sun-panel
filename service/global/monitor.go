package global

import "sun-panel/lib/monitor"

type ModelSystemMonitor struct {
	CPUInfo           monitor.CPUInfo             `json:"cpuInfo"`
	DiskInfo          []monitor.DiskInfo          `json:"diskInfo"`
	NetIOCountersInfo []monitor.NetIOCountersInfo `json:"netIOCountersInfo"`
	MemoryInfo        monitor.MemoryInfo          `json:"memoryInfo"`
}
