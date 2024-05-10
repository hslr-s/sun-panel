package biz

import (
	"context"
	"encoding/json"
	"io"
	"strings"
	"sun-panel/api/api_v1/common/apiData/dockerApiStructs"
	"sun-panel/lib/docker"

	"github.com/docker/docker/api/types"
)

// docker 官网文档
// https://docs.docker.com/engine/api/

type DockerType struct {
}

func (g *DockerType) GetContainerInfo(containerId string) (dockerApiStructs.ContainerInfo, error) {
	var data dockerApiStructs.ContainerInfo

	client, err := docker.NewDockerClient()
	if err != nil {
		return data, err
	}
	defer client.Close()

	ctx := context.Background()
	container, err := client.ContainerInspect(ctx, containerId)
	if err != nil {
		return data, err
	}
	data.Status = container.State.Status
	data.ID = container.ID
	data.Name = container.Name
	data.Image = container.Image
	data.CreateTime = container.Created
	return data, nil
}

func (g *DockerType) GetContainerInspect(containerId string) (dockerApiStructs.ContainerInfo, error) {
	var data dockerApiStructs.ContainerInfo

	client, err := docker.NewDockerClient()
	if err != nil {
		return data, err
	}
	defer client.Close()

	ctx := context.Background()

	container, err := client.ContainerInspect(ctx, containerId)
	if err != nil {
		return data, err
	}
	data.Status = container.State.Status
	data.ID = container.ID
	data.Name = container.Name
	return data, nil
}

func (g *DockerType) GetContainerState(containerId string) (dockerApiStructs.ContainerStats, error) {
	var data dockerApiStructs.ContainerStats

	client, err := docker.NewDockerClient()
	if err != nil {
		return data, err
	}
	defer client.Close()

	ctx := context.Background()

	res, err := client.ContainerStats(ctx, containerId, false)
	if err != nil {
		return data, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		res.Body.Close()
		return data, err
	}

	res.Body.Close()
	var stats *types.StatsJSON
	if err := json.Unmarshal(body, &stats); err != nil {
		return data, err
	}

	data.CPUPercent = calculateCPUPercentUnix(stats)
	data.IORead, data.IOWrite = calculateBlockIO(stats.BlkioStats)
	data.Memory = float64(stats.MemoryStats.Usage) / 1024 / 1024
	if cache, ok := stats.MemoryStats.Stats["cache"]; ok {
		data.Cache = float64(cache) / 1024 / 1024
	}
	data.Memory = data.Memory - data.Cache
	data.NetworkRX, data.NetworkTX = calculateNetwork(stats.Networks)
	data.ShotTime = stats.Read
	return data, nil
}

func calculateCPUPercentUnix(stats *types.StatsJSON) float64 {
	cpuPercent := 0.0
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage) - float64(stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage) - float64(stats.PreCPUStats.SystemUsage)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * 100.0
		if len(stats.CPUStats.CPUUsage.PercpuUsage) != 0 {
			cpuPercent = cpuPercent * float64(len(stats.CPUStats.CPUUsage.PercpuUsage))
		}
	}
	return cpuPercent
}
func calculateMemPercentUnix(memStats types.MemoryStats) float64 {
	memPercent := 0.0
	memUsage := float64(memStats.Usage - memStats.Stats["cache"])
	memLimit := float64(memStats.Limit)
	if memUsage > 0.0 && memLimit > 0.0 {
		memPercent = (memUsage / memLimit) * 100.0
	}
	return memPercent
}
func calculateBlockIO(blkio types.BlkioStats) (blkRead float64, blkWrite float64) {
	for _, bioEntry := range blkio.IoServiceBytesRecursive {
		switch strings.ToLower(bioEntry.Op) {
		case "read":
			blkRead = (blkRead + float64(bioEntry.Value)) / 1024 / 1024
		case "write":
			blkWrite = (blkWrite + float64(bioEntry.Value)) / 1024 / 1024
		}
	}
	return
}
func calculateNetwork(network map[string]types.NetworkStats) (float64, float64) {
	var rx, tx float64

	for _, v := range network {
		rx += float64(v.RxBytes) / 1024
		tx += float64(v.TxBytes) / 1024
	}
	return rx, tx
}
