package openApi

import (
	v1 "sun-panel/router/openApi/v1"

	"github.com/gin-gonic/gin"
)

func Init(routerGroup *gin.RouterGroup) {
	openApi := routerGroup.Group("openapi")
	v1.Init(openApi)
}
