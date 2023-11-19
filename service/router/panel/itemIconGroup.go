package panel

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitItemIconGroup(router *gin.RouterGroup) {
	itemIconGroup := api_v1.ApiGroupApp.ApiPanel.ItemIconGroup
	r := router.Group("", middleware.LoginInterceptor)
	{
		r.POST("/panel/itemIconGroup/edit", itemIconGroup.Edit)
		r.POST("/panel/itemIconGroup/getList", itemIconGroup.GetList)
		r.POST("/panel/itemIconGroup/deletes", itemIconGroup.Deletes)
	}
}
