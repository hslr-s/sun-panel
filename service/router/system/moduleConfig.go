package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitModuleConfigRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.ModuleConfigApi
	r := router.Group("", middleware.LoginInterceptorDev)
	r.POST("/system/moduleConfig/getByName", api.GetByName)
	r.POST("/system/moduleConfig/save", api.Save)

}
