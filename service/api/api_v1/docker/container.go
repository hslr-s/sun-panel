package docker

import (
	"context"
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
