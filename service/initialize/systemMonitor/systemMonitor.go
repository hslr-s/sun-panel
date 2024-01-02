package systemMonitor

import (
	"sun-panel/global"
	"sun-panel/lib/cache"
	"sun-panel/lib/monitor"
	"time"
)

func Start(cacher cache.Cacher[global.ModelSystemMonitor], interval time.Duration) {
	go func() {

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				go func() {
					GetInfo()
				}()
			}
		}

	}()

}

func GetInfo() global.ModelSystemMonitor {

	var modelSystemMonitor global.ModelSystemMonitor

	if cpuInfo, err := monitor.GetCPUInfo(); err == nil {
		modelSystemMonitor.CPUInfo = cpuInfo
	}

	if v, err := monitor.GetDiskInfo(); err == nil {
		modelSystemMonitor.DiskInfo = v
	}

	if v, err := monitor.GetNetIOCountersInfo(); err == nil {
		modelSystemMonitor.NetIOCountersInfo = v
	}

	if v, err := monitor.GetMemoryInfo(); err == nil {
		modelSystemMonitor.MemoryInfo = v
	}

	return modelSystemMonitor
}
