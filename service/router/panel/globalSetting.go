package panel

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitGetSiteStting(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiPanel.GlobalSettingApi
	r := router.Group("", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("/panel/globalSetting/saveSiteStting", api.SaveSiteStting)
		r.POST("/panel/globalSetting/uploadFaviconImage", api.UploadFaviconImage)
		r.POST("/panel/globalSetting/getCustomStylePath", api.GetCustomStylePath)
		r.POST("/panel/globalSetting/saveCustomJsAndCssCode", api.SaveCustomJsAndCssCode)
		r.POST("/panel/globalSetting/getCustomJsAndCssCode", api.GetCustomJsAndCssCode)
	}

	// 公开模式
	rPublic := router.Group("")
	{
		rPublic.POST("/panel/globalSetting/getSiteStting", api.GetSiteStting)
	}
}
