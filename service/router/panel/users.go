package panel

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitUsersRouter(router *gin.RouterGroup) {
	userApi := api_v1.ApiGroupApp.ApiPanel.UsersApi

	// r := router.Group("", middleware.LoginInterceptor)
	rAdmin := router.Group("", middleware.AdminInterceptor)
	{
		rAdmin.POST("panel/users/create", userApi.Create)
		rAdmin.POST("panel/users/update", userApi.Update)
		rAdmin.POST("panel/users/getList", userApi.GetList)
		rAdmin.POST("panel/users/deletes", userApi.Deletes)
	}
}
