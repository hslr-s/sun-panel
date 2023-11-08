package panel

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitUsersRouter(router *gin.RouterGroup) {
	userApi := api_v1.ApiGroupApp.ApiPanel.UsersApi

	r := router.Group("", middleware.LoginInterceptor)
	{
		r.POST("panel/users/create", userApi.Create)
		r.POST("panel/users/update", userApi.Update)
		r.POST("panel/users/getList", userApi.GetList)
		r.POST("panel/users/deletes", userApi.Deletes)
	}
}
