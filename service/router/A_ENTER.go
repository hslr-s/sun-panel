package router

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sun-panel/biz"
	"sun-panel/global"

	// "sun-panel/router/admin"

	"sun-panel/router/docker"
	"sun-panel/router/openApi"
	"sun-panel/router/openness"
	"sun-panel/router/panel"
	"sun-panel/router/system"

	"github.com/gin-gonic/gin"
)

// 初始化总路由
func InitRouters() {

	addr := global.Config.GetValueStringOrDefault("base", "http_address")
	httpPort := global.Config.GetValueInt("base", "http_port")
	httpsPort := global.Config.GetValueIntOrDefault("base", "https_port", 3003)

	certPath := global.Config.GetValueString("cert", "cert_path")
	keyPath := global.Config.GetValueString("cert", "key_path")

	httpRouter := gin.Default()
	httpsRouter := gin.Default()

	httpRouter = InitRouter(httpRouter)
	httpsRouter = InitRouter(httpsRouter)

	httpAddress := fmt.Sprintf("%s:%d", addr, httpPort)
	httpsAddress := fmt.Sprintf("%s:%d", addr, httpsPort)

	global.Logger.Info("Sun-Panel is Started.")

	global.Logger.Info("HTTP: ", httpAddress)

	// 开启HTTPS服务
	if certPath != "" {
		go func() {
			global.Logger.Info("HTTPS: ", httpsAddress)
			err := httpsRouter.RunTLS(fmt.Sprintf("%s:%d", addr, httpsPort), certPath, keyPath)
			log.Fatalln(err)
		}()
	}

	err := httpRouter.Run(fmt.Sprintf("%s:%d", addr, httpPort))
	log.Fatalln(err)
}

func InitRouter(router *gin.Engine) *gin.Engine {
	rootRouter := router.Group("/")
	routerGroup := rootRouter.Group("api")

	// 接口
	system.Init(routerGroup)
	panel.Init(routerGroup)
	docker.Init(routerGroup)
	openness.Init(routerGroup)

	// 开放API
	openApi.Init(rootRouter)

	// WEB文件服务
	{
		webPath := global.Config.GetValueStringOrDefault("base", "web_path")

		// 要跳过的目录列表
		skipDirs := map[string]bool{
			"custom":  true,
			"uploads": true,
		}

		// 要跳过的文件列表
		skipFiles := map[string]bool{
			// "favicon.ico": true,
		}

		// 动态绑定web目录的路由
		if err := bindWebRouters(router, webPath, skipDirs, skipFiles); err != nil {
			log.Fatalf("Error reading directory '%s': %v\n", webPath, err)
		}

		// 后端无法匹配的路由，将其重定向到前端
		router.NoRoute(func(c *gin.Context) {
			// 兼容前端history模式，后端路由匹配不到将其重定向到前端index.html
			// c.File(webPath + "/index.html")

			// 读取文件并替换模版变量
			content, err := os.ReadFile(webPath + "/index.html")
			if err != nil {
				c.String(500, "Failed to load file: %v", err)
				return
			}

			globalSettingSiteSetting := biz.GlobalSettingSiteSetting{}
			biz.GlobalSetting.GetSiteStting(&globalSettingSiteSetting)

			// 设置默认值
			if globalSettingSiteSetting.FaviconUrl == "" || globalSettingSiteSetting.FaviconUrl == "favicon.ico" || globalSettingSiteSetting.FaviconUrl == "./favicon.ico" {
				globalSettingSiteSetting.FaviconUrl = "/favicon.ico"
			}

			if globalSettingSiteSetting.Title == "" {
				globalSettingSiteSetting.Title = "Sun-Panel"
			}

			// 替换内容中的占位符 ( 模版变量格式:{{__VAR__}} )
			htmlContent := string(content)
			htmlContent = strings.ReplaceAll(htmlContent, "{{__title__}}", globalSettingSiteSetting.Title)
			htmlContent = strings.ReplaceAll(htmlContent, "{{__faviconSrc__}}", globalSettingSiteSetting.FaviconUrl)

			// 返回修改后的 HTML
			c.Data(200, "text/html; charset=utf-8", []byte(htmlContent))
		})

		// 其他特殊要绑定的路由

		// 绑定上传目录的路由
		sourcePath := global.Config.GetValueString("base", "source_path")
		router.Static("/uploads", sourcePath)

		// 绑定自定义风格的路由
		customPath := global.Config.GetValueString("base", "custom_style_path")
		// 兼容旧版 2024-02-06 开发
		if customPath == "" {
			router.Static("/custom", webPath+"/custom")
		} else {
			router.Static("/custom", customPath)
		}

	}

	return router
}

// 动态绑定web目录的路由
func bindWebRouters(router *gin.Engine, webPath string, skipDirs, skipFiles map[string]bool) error {

	// 遍历 webPath 目录的一级内容，动态绑定到路由
	entries, err := os.ReadDir(webPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		// 获取相对路径
		relativePath := "/" + entry.Name()
		fullPath := filepath.Join(webPath, entry.Name())

		// 检查是否在跳过的目录列表中
		if entry.IsDir() && skipDirs[entry.Name()] {
			continue
		}

		// 检查是否在跳过的文件列表中
		if !entry.IsDir() && skipFiles[entry.Name()] {
			continue
		}

		if entry.IsDir() {
			// global.Logger.Debug("动态绑定目录", relativePath, fullPath)
			// 如果是目录，使用 Static 绑定目录
			router.Static(relativePath, fullPath)
		} else {
			// global.Logger.Debug("动态绑定文件", relativePath, fullPath)
			// 如果是文件，使用 StaticFile 绑定文件
			router.StaticFile(relativePath, fullPath)
		}
	}

	return nil
}
