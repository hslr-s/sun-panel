package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.UserApi
	r := router.Group("", middleware.LoginInterceptor)
	r.POST("/user/getInfo", api.GetInfo)
	r.POST("/user/updatePasssword", api.UpdatePasssword)
	r.POST("/user/updateInfo", api.UpdateInfo)
	r.POST("/user/getReferralCode", api.GetReferralCode)
}
