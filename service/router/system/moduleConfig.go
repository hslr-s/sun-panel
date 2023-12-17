package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitModuleConfigRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.ModuleConfigApi
	r := router.Group("", middleware.LoginInterceptor)
	r.POST("/system/moduleConfig/save", api.Save)

	// 公开模式
	rPublic := router.Group("", middleware.PublicModeInterceptor)
	{
		rPublic.POST("/system/moduleConfig/getByName", api.GetByName)
	}
}
