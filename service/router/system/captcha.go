package system

import (
	"sun-panel/api/api_v1"

	"github.com/gin-gonic/gin"
)

func InitCaptchaRouter(router *gin.RouterGroup) {
	captchaApi := api_v1.ApiGroupApp.ApiSystem.CaptchaApi
	r := router.Group("captcha")
	r.GET("getImage", captchaApi.GetImage)
	r.GET("getImage/:width/:height", captchaApi.GetImage)
	r.GET("getImageByCaptchaId/:captchaId", captchaApi.GetImageByCaptchaId)
	r.GET("getImageByCaptchaId/:captchaId/:width/:height", captchaApi.GetImageByCaptchaId)
	// r.POST("/captach/check", captchaApi.CheckVCode)
}
