package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitOpenApiRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.OpenApiApi
	r := router.Group("", middleware.LoginInterceptor)

	r.POST("/openApi/refreshToken", api.RefreshToken)
	r.POST("/openApi/deleteToken", api.DeleteToken)
	r.POST("/openApi/getTokenStatusInfo", api.GetTokenStatusInfo)

}
