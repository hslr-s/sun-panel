package panel

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitItemIcon(router *gin.RouterGroup) {
	itemIcon := api_v1.ApiGroupApp.ApiPanel.ItemIcon
	r := router.Group("", middleware.LoginInterceptor)
	{
		r.POST("/panel/itemIcon/edit", itemIcon.Edit)
		r.POST("/panel/itemIcon/deletes", itemIcon.Deletes)
		r.POST("/panel/itemIcon/saveSort", itemIcon.SaveSort)
		r.POST("/panel/itemIcon/addMultiple", itemIcon.AddMultiple)
		r.POST("/panel/itemIcon/getSiteFavicon", itemIcon.GetSiteFavicon)
	}

	// 公开模式
	rPublic := router.Group("", middleware.PublicModeInterceptor)
	{
		rPublic.POST("/panel/itemIcon/getListByGroupId", itemIcon.GetListByGroupId)
	}
}
