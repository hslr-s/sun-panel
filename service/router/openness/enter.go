package openness

import "github.com/gin-gonic/gin"

func Init(routerGroup *gin.RouterGroup) {
	opennessGroup := routerGroup.Group("openness")
	InitOpenness(opennessGroup)
}
