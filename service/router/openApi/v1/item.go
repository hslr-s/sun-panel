package v1

import (
	"sun-panel/api/openAPI_v1"
	"sun-panel/api/openAPI_v1/common/middleware"

	"github.com/gin-gonic/gin"
)

func InitItemRouter(routerGroup *gin.RouterGroup) {
	api := openAPI_v1.OpenApiGroupApp.PanelApiGroup.ItemApi
	r := routerGroup.Group("/item", middleware.LoginInterceptor)
	{
		r.POST("create", api.Create)
		r.POST("update", api.Update)
		r.POST("getInfoByOnlyName", api.GetInfoByOnlyName)
	}
}
