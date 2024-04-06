package router

import (
	"sun-panel/biz"
	"sun-panel/global"

	// "sun-panel/router/admin"
	"sun-panel/router/openApi"
	"sun-panel/router/openness"
	"sun-panel/router/panel"
	"sun-panel/router/system"

	"github.com/gin-gonic/gin"
)

// 初始化总路由
func InitRouters(addr string) error {
	router := gin.Default()
	rootRouter := router.Group("/")
	routerGroup := rootRouter.Group("api")

	// 接口
	system.Init(routerGroup)
	panel.Init(routerGroup)
	openness.Init(routerGroup)
	openApi.Init(routerGroup)

	// WEB文件服务
	{
		webPath := global.Config.GetValueStringOrDefault("base", "web_path")
		// 上传
		sourcePath := global.Config.GetValueString("base", "source_path")
		router.Static("/uploads", sourcePath)
		// 自定义风格文件夹
		customPath := global.Config.GetValueString("base", "custom_style_path")
		// 兼容旧版
		if customPath == "" {
			router.Static("/custom", webPath+"/custom")
		} else {
			router.Static("/custom", customPath)
		}
		router.StaticFile("/", webPath+"/index.html")
		router.Static("/assets", webPath+"/assets")
		// router.StaticFile("/favicon.ico", webPath+"/favicon.ico")
		// router.StaticFile("/favicon.svg", webPath+"/favicon.svg")

		// 动态返回favicon
		router.Any("/favicon.ico", func(c *gin.Context) {
			setting := biz.GlobalSettingSiteSetting{}
			biz.GlobalSetting.GetSiteStting(&setting)
			if setting.FaviconUrl == "" || setting.FaviconUrl == "favicon.ico" || setting.FaviconUrl == "./favicon.ico" {
				setting.FaviconUrl = webPath + "/favicon.ico"
			}
			c.File(setting.FaviconUrl)
		})
	}

	global.Logger.Info("Sun-Panel is Started.  Listening and serving HTTP on 0.0.0.0", addr)
	return router.Run(addr)
}
