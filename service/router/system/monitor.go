package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitMonitorRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.MonitorApi
	r := router.Group("", middleware.LoginInterceptor)
	r.POST("/system/monitor/getDiskMountpoints", api.GetDiskMountpoints)

	// 公开模式
	rPublic := router.Group("", middleware.PublicModeInterceptor)
	{
		rPublic.POST("/system/monitor/getAll", api.GetAll)
		rPublic.POST("/system/monitor/getCpuState", api.GetCpuState)
		rPublic.POST("/system/monitor/getDiskStateByPath", api.GetDiskStateByPath)
		rPublic.POST("/system/monitor/getMemonyState", api.GetMemonyState)
	}
}
