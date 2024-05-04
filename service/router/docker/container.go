package docker

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitContainer(router *gin.RouterGroup) {
	apiDocker := api_v1.ApiGroupApp.ApiDocker
	r := router.Group("", middleware.LoginInterceptor)

	// r := router.Group("", middleware.LoginInterceptor)
	{
		r.POST("/docker/container/getList", apiDocker.Container.GetList)
		r.POST("/docker/container/execSwitchActionByContainerID", apiDocker.Container.ExecSwitchActionByContainerID)
		r.POST("/docker/container/getContainerStatus", apiDocker.Container.GetContainerStatus)

		// r.POST("/docker/itemIcon/getListByGroupId", itemIcon.GetListByGroupId)
		// r.POST("/docker/itemIcon/deletes", itemIcon.Deletes)
	}
}
