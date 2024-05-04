package docker

import (
	"context"
	"encoding/json"
	"io"
	"strings"
	"sun-panel/api/api_v1/common/apiData/dockerApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/lib/cmn"
	"sun-panel/lib/docker"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Container struct {
}

func (a *Container) GetList(c *gin.Context) {
	// req := commonApiStructs.RequestPage{}
	// if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
	// 	apiReturn.ErrorParamFomat(c, err.Error())
	// 	return
	// }

	var (
		client *client.Client
		err    error
	)
	client, err = docker.NewDockerClient()
	if err != nil {
		apiReturn.ErrorDockerClientFail(c, err.Error())
		return
	}

	options := types.ContainerListOptions{
		All:  true,
		Size: true,
	}
	if list, err := client.ContainerList(context.Background(), options); err != nil {
		apiReturn.ErrorDocker(c, err.Error())
		return
	} else {
		resp := []dockerApiStructs.ContainerInfoResp{}
		for _, v := range list {
			ports := []dockerApiStructs.ContainerPort{}
			for _, vPort := range v.Ports {
				ports = append(ports, dockerApiStructs.ContainerPort{
					IP:          vPort.IP,
					PrivatePort: vPort.PrivatePort,
					PublicPort:  vPort.PublicPort,
					Type:        dockerApiStructs.PostType(vPort.Type),
				})
			}
			resp = append(resp, dockerApiStructs.ContainerInfoResp{
				ID:         v.ID,
				Name:       v.Names[0][1:],
				Image:      v.Image,
				State:      v.State,
				Status:     v.Status,
				CreateTime: time.Unix(v.Created, 0).Format(cmn.TimeFormatMode1),
				Ports:      ports,
			})
		}
		apiReturn.SuccessListData(c, resp, 0)
	}
}

func (a *Container) ExecSwitchActionByContainerID(c *gin.Context) {
	req := dockerApiStructs.ExecSwitchActionByContainerIDReq{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	var (
		client *client.Client
		err    error
	)
	client, err = docker.NewDockerClient()
	if err != nil {
		apiReturn.ErrorDockerClientFail(c, err.Error())
		return
	}
	defer client.Close()

	ctx := context.Background()

	if req.Action {
		startOptions := types.ContainerStartOptions{}
		err = client.ContainerStart(ctx, req.ID, startOptions)
	} else {
		startOptions := container.StopOptions{}
		err = client.ContainerStop(ctx, req.ID, startOptions)
	}

	if err != nil {
		apiReturn.ErrorDocker(c, err.Error())
		return
	}
	apiReturn.Success(c)

}

// 获取容器状态
func (a *Container) GetContainerStatus(c *gin.Context) {
	req := dockerApiStructs.ExecSwitchActionByContainerIDReq{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	var (
		client *client.Client
		err    error
	)
	client, err = docker.NewDockerClient()
	if err != nil {
		apiReturn.ErrorDockerClientFail(c, err.Error())
		return
	}
	defer client.Close()

	ctx := context.Background()

	res, err := client.ContainerStats(ctx, req.ID, false)
	if err != nil {
		apiReturn.ErrorDocker(c, err.Error())
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		res.Body.Close()
		apiReturn.ErrorDocker(c, err.Error())
		return
	}
	res.Body.Close()
	var stats *types.StatsJSON
	if err := json.Unmarshal(body, &stats); err != nil {
		apiReturn.ErrorDocker(c, err.Error())
		return
	}
	var data dockerApiStructs.ContainerStatsRespItem
	data.CPUPercent = calculateCPUPercentUnix(stats)
	data.IORead, data.IOWrite = calculateBlockIO(stats.BlkioStats)
	data.Memory = float64(stats.MemoryStats.Usage) / 1024 / 1024
	if cache, ok := stats.MemoryStats.Stats["cache"]; ok {
		data.Cache = float64(cache) / 1024 / 1024
	}
	data.Memory = data.Memory - data.Cache
	data.NetworkRX, data.NetworkTX = calculateNetwork(stats.Networks)
	data.ShotTime = stats.Read

	apiReturn.SuccessData(c, data)

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
