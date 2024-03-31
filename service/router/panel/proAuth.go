package panel

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitProAuthRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiPanel.ProAuthApi

	r := router.Group("", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("/admin/proAuth/login", api.Login)
		r.POST("/admin/proAuth/getStatus", api.GetStatus)
	}

}
