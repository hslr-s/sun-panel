package monitor

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type CPUInfo struct {
	CoreCount int32     `json:"coreCount"`
	CPUNum    int       `json:"cpuNum"`
	Model     string    `json:"model"`
	Usages    []float64 `json:"usages"`
}

type DiskInfo struct {
	Mountpoint  string  `json:"mountpoint"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"usedPercent"`
}

type NetIOCountersInfo struct {
	BytesSent uint64 `json:"bytesSent"`
	BytesRecv uint64 `json:"bytesRecv"`
	Name      string `json:"name"`
}

type MemoryInfo struct {
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

// 获取CPU信息
func GetCPUInfo() (CPUInfo, error) {
	cpuInfoRes := CPUInfo{}
	cpuInfo, err := cpu.Info()

	if err == nil && len(cpuInfo) > 0 {
		cpuInfoRes.CoreCount = cpuInfo[0].Cores
		cpuInfoRes.Model = cpuInfo[0].ModelName
	}
	numCPU, _ := cpu.Counts(true)
	cpuInfoRes.CPUNum = numCPU
	cpuPercentages, err := cpu.Percent(time.Second, true)
	cpuInfoRes.Usages = cpuPercentages

	return cpuInfoRes, err
}

// 获取内存信息 单位：MB
func GetMemoryInfo() (MemoryInfo, error) {
	memoryInfo := MemoryInfo{}
	// 获取内存信息
	memInfo, err := mem.VirtualMemory()
	if err == nil {
		memoryInfo.Free = memInfo.Free
		memoryInfo.Total = memInfo.Total
		memoryInfo.Used = memInfo.Used
		memoryInfo.UsedPercent = memInfo.UsedPercent
	}

	return memoryInfo, err
}

// 获取每个磁盘分区使用情况
func GetDiskInfo() ([]DiskInfo, error) {
	disks := []DiskInfo{}
	// 获取所有磁盘分区的信息
	partitions, err := disk.Partitions(true)
	if err != nil {
		return disks, err
	}

	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			// fmt.Printf("Error getting disk usage for %s: %v\n", partition.Mountpoint, err)
			continue
		}

		disks = append(disks, DiskInfo{
			Mountpoint:  partition.Mountpoint,
			Total:       usage.Total / 1024 / 1024,
			Used:        usage.Used / 1024 / 1024,
			Free:        usage.Free / 1024 / 1024,
			UsedPercent: usage.UsedPercent,
		})
	}

	return disks, nil
}

func GetDiskMountpoints() ([]disk.PartitionStat, error) {
	return disk.Partitions(true)
}

func GetDiskInfoByPath(path string) (*DiskInfo, error) {
	diskInfo := DiskInfo{}
	usage, err := disk.Usage(path)
	if err != nil {
		return nil, err
	}
	diskInfo.Free = usage.Free
	diskInfo.Mountpoint = usage.Path
	diskInfo.Total = usage.Total
	diskInfo.Used = usage.Used
	diskInfo.UsedPercent = usage.UsedPercent
	return &diskInfo, nil
}

// 获取网络统计信息
func GetNetIOCountersInfo() ([]NetIOCountersInfo, error) {
	netInfo := []NetIOCountersInfo{}
	netStats, err := net.IOCounters(true)
	if err == nil {
		for _, netStat := range netStats {
			netInfo = append(netInfo, NetIOCountersInfo{
				BytesRecv: netStat.BytesRecv,
				BytesSent: netStat.BytesSent,
				Name:      netStat.Name,
			})

		}
	}
	return netInfo, err
}

// func GetCountDiskInfo() {
// 	// 获取所有磁盘的总使用情况
// 	allUsage, err := disk.Usage("/")
// 	if err != nil {
// 		fmt.Printf("Error getting total disk usage: %v\n", err)
// 		return
// 	}

// 	// 打印所有磁盘的总使用情况
// 	fmt.Println("Total Disk Usage:")
// 	fmt.Printf("Total: %d MB\n", allUsage.Total/1024/1024)
// 	fmt.Printf("Used: %d MB\n", allUsage.Used/1024/1024)
// 	fmt.Printf("Free: %d MB\n", allUsage.Free/1024/1024)
// 	fmt.Printf("Usage: %.2f%%\n", allUsage.UsedPercent)
// }
