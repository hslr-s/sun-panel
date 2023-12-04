package panel

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitUsersRouter(router *gin.RouterGroup) {
	userApi := api_v1.ApiGroupApp.ApiPanel.UsersApi

	rAdmin := router.Group("", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		rAdmin.POST("panel/users/create", userApi.Create)
		rAdmin.POST("panel/users/update", userApi.Update)
		rAdmin.POST("panel/users/getList", userApi.GetList)
		rAdmin.POST("panel/users/deletes", userApi.Deletes)
		rAdmin.POST("panel/users/getPublicVisitUser", userApi.GetPublicVisitUser)
		rAdmin.POST("panel/users/setPublicVisitUser", userApi.SetPublicVisitUser)
	}
}
