package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitAbout(router *gin.RouterGroup) {
	about := api_v1.ApiGroupApp.ApiSystem.About
	{
		router.POST("about", about.Get)

		// 需要登录可以调用的接口
		lr := router.Group("about/", middleware.LoginInterceptor)
		lr.POST("version/checkNewVersion", about.CheckNewVersion)
		lr.POST("version/setAutoCheckNewVersion", about.SetAutoCheckNewVersion)
		lr.POST("version/getAutoCheckNewVersion", about.GetAutoCheckNewVersion)
	}
}
