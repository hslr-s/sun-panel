package docker

import (
	"context"
	"sun-panel/api/api_v1/common/apiData/dockerApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/lib/cmn"
	"sun-panel/lib/docker"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
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
			resp = append(resp, dockerApiStructs.ContainerInfoResp{
				Name:       v.Names[0][1:],
				Image:      v.Image,
				State:      v.State,
				CreateTime: time.Unix(v.Created, 0).Format(cmn.TimeFormatMode1),
			})
		}
		apiReturn.SuccessListData(c, resp, 0)
	}
}
