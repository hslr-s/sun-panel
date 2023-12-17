package system

import "github.com/gin-gonic/gin"

func Init(routerGroup *gin.RouterGroup) {
	InitAbout(routerGroup)
	InitLogin(routerGroup)
	InitUserRouter(routerGroup)
	InitFileRouter(routerGroup)
	InitCaptchaRouter(routerGroup)
	InitRegister(routerGroup)
	InitNoticeRouter(routerGroup)
	InitModuleConfigRouter(routerGroup)
}
