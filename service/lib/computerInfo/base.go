package computerInfo

import (
	"log"
	"time"

	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/v3/cpu"
	"gitlab.com/tingshuo/go-diskstate/diskstate"
)

type Storage struct {
	Name       string
	FileSystem string
	Total      uint64
	Free       uint64
}

type storageInfo struct {
	Name       string
	Size       uint64
	FreeSpace  uint64
	FileSystem string
	Used       uint64
}

// func GetStorageInfo() {
// 	var storageinfo []storageInfo
// 	var loaclStorages []Storage
// 	err := wmi.Query("Select * from Win32_LogicalDisk", &storageinfo)
// 	if err != nil {
// 		return
// 	}

// 	for _, storage := range storageinfo {
// 		info := Storage{
// 			Name:       storage.Name,
// 			FileSystem: storage.FileSystem,
// 			Total:      storage.Size / 1024 / 1024 / 1024,
// 			Free:       storage.FreeSpace / 1024 / 1024 / 1024,
// 		}
// 		if info.Total >= 1 {
// 			fmt.Printf("%s总大小%dG，可用%dG\n", info.Name, info.Total, info.Free)
// 			loaclStorages = append(loaclStorages, info)
// 		}
// 	}
// 	fmt.Printf("localStorages:= %v\n", loaclStorages)
// }

func GetCurrentStorageInfo(path string) storageInfo {
	state := diskstate.DiskUsage(path)
	info := storageInfo{}
	info.Size = uint64(state.All / diskstate.B)
	info.FreeSpace = uint64(state.Free / diskstate.B)
	info.Used = uint64(state.Used / diskstate.B)

	// fmt.Printf("All=%dM, Free=%dM, Available=%dM, Used=%dM, Usage=%d%%",
	// 	state.All/diskstate.B, state.Free/diskstate.MB, state.Available/diskstate.MB, state.Used/diskstate.MB, 100*state.Used/state.All)
	return info
}

type ComputerMonitor struct {
	CPU float64 `json:"cpu"`
	Mem float64 `json:"mem"`
}

// GetCPUPercent 获取CPU使用率
func GetCPUPercent() float64 {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Fatalln(err.Error())
		return -1
	}
	return percent[0]
}

// GetMemPercent 获取内存使用率
func GetMemPercent() float64 {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		log.Fatalln(err.Error())
		return -1
	}
	return memInfo.UsedPercent
}

func GetCpuMem() ComputerMonitor {
	var res ComputerMonitor
	res.CPU = GetCPUPercent()
	res.Mem = GetMemPercent()
	return res
}
