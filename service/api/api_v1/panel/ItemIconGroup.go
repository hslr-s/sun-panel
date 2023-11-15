package panel

import (
	"sun-panel/api/api_v1/common/apiData/commonApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type ItemIconGroup struct {
}

func (a *ItemIconGroup) Edit(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	req := models.ItemIconGroup{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	req.UserId = userInfo.ID

	if req.ID != 0 {
		// 修改
		global.Db.Model(&models.ItemIconGroup{}).
			Select("IconJson", "Icon", "Title", "Url", "LanUrl", "Description", "OpenMethod", "Sort", "GroupId", "UserId").
			Where("id=?", req.ID).Updates(&req)
	} else {
		// 创建
		global.Db.Create(&req)
	}

	apiReturn.SuccessData(c, req)
}

func (a *ItemIconGroup) GetList(c *gin.Context) {

	userInfo, _ := base.GetCurrentUserInfo(c)
	groups := []models.ItemIconGroup{}

	if err := global.Db.Order("sort ,created_at").Where("user_id=?", userInfo.ID).Find(&groups, "user_id=?", 1, userInfo.ID).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	// 判断分组是否为空，为空将自动创建默认分组
	if len(groups) == 0 {
		defaultGroup := models.ItemIconGroup{
			Title:  "APP",
			UserId: userInfo.ID,
			Icon:   "material-symbols:folder-outline"}
		if err := global.Db.Create(&defaultGroup).Error; err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
			return
		}

		// 并将当前账号下所有无分组的图标更新到当前组
		if err := global.Db.Model(&models.ItemIcon{}).Where("user_id=?", userInfo.ID).Update("item_icon_group_id", defaultGroup.ID).Error; err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
			return
		}

		groups = append(groups, defaultGroup)
	}

	apiReturn.SuccessListData(c, groups, 0)
}

func (a *ItemIconGroup) Deletes(c *gin.Context) {
	req := commonApiStructs.RequestDeleteIds[uint]{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)
	if err := global.Db.Delete(&models.ItemIconGroup{}, "id in ? AND user_id=?", req.Ids, userInfo.ID).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}
