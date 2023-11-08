package openness

import (
	"sun-panel/api/api_v1"

	"github.com/gin-gonic/gin"
)

func InitOpenness(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiOpen.Openness
	{
		router.GET("loginConfig", api.LoginConfig)
		router.GET("getDisclaimer", api.GetDisclaimer)
		router.GET("getAboutDescription", api.GetAboutDescription)
	}
}
