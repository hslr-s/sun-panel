package v1

import (
	"sun-panel/api/openAPI_v1"
	"sun-panel/api/openAPI_v1/common/middleware"

	"github.com/gin-gonic/gin"
)

func InitItemGroupRouter(routerGroup *gin.RouterGroup) {
	api := openAPI_v1.OpenApiGroupApp.PanelApiGroup.ItemGroupApi
	r := routerGroup.Group("/itemGroup", middleware.LoginInterceptor)
	{
		r.POST("create", api.Create)
		// r.POST("update", api.Update)
		r.POST("getInfo", api.GetInfo)
		r.POST("getList", api.GetList)
	}
}
