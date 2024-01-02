package global

import "sun-panel/lib/monitor"

type ModelSystemMonitor struct {
	CPUInfo           monitor.CPUInfo
	DiskInfo          []monitor.DiskInfo
	NetIOCountersInfo []monitor.NetIOCountersInfo
	MemoryInfo        monitor.MemoryInfo
	// NetIOCountersInfo monitor.NetIOCountersInfo

}
