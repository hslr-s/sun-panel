package v1

import "github.com/gin-gonic/gin"

func Init(routerGroup *gin.RouterGroup) {
	v1Group := routerGroup.Group("v1")
	InitItemRouter(v1Group)
	InitItemGroupRouter(v1Group)
	InitVersionRouter(v1Group)
}
