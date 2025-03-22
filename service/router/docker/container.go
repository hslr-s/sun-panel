package docker

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitContainer(router *gin.RouterGroup) {
	apiDocker := api_v1.ApiGroupApp.ApiDocker

	// 登录模式（可以获取docker列表，可以控制容器）
	r := router.Group("", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("/docker/container/getList", apiDocker.Container.GetList)
		r.POST("/docker/container/execSwitchActionByContainerID", apiDocker.Container.ExecSwitchActionByContainerID)
		r.POST("/docker/container/getContainerIDByName", apiDocker.Container.GetContainerIDByName)
	}

	// 公开模式（仅可以获取容器的状态）
	pm := router.Group("", middleware.PublicModeInterceptor)
	{
		pm.POST("/docker/container/getContainerStatus", apiDocker.Container.GetContainerStatus)
		pm.POST("/docker/container/getContainersStats", apiDocker.Container.GetContainersState)
	}

}
