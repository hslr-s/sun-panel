package docker

import "github.com/gin-gonic/gin"

func Init(routerGroup *gin.RouterGroup) {
	InitContainer(routerGroup)
}
