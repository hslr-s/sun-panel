package v1

import (
	"sun-panel/api/openAPI_v1"
	"sun-panel/api/openAPI_v1/common/middleware"

	"github.com/gin-gonic/gin"
)

func InitVersionRouter(routerGroup *gin.RouterGroup) {
	api := openAPI_v1.OpenApiGroupApp.PanelApiGroup.VersionApi
	r := routerGroup.Group("/version", middleware.LoginInterceptor)
	{
		r.Any("get", api.GetInfo)
		r.Any("", api.GetInfo)
	}
}
