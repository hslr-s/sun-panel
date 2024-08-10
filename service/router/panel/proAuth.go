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
		r.POST("/panel/proAuth/login", api.Login)
		r.POST("/panel/proAuth/getStatus", api.GetStatus)
		r.POST("/panel/proAuth/logout", api.Logout)
		r.POST("/panel/proAuth/refreshStatus", api.RefreshStatus)
		r.POST("/panel/proAuth/setHideProBadgeStatus", api.SetHideProBadgeStatus)
	}

}
