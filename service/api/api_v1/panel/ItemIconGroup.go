package panel

import (
	"errors"
	"math"
	"sun-panel/api/api_v1/common/apiData/commonApiStructs"
	"sun-panel/api/api_v1/common/apiData/panelApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type ItemIconGroup struct {
}

func (a *ItemIconGroup) Edit(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	req := panelApiStructs.ItemGroupEditReq{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	req.UserId = userInfo.ID

	saveData := models.ItemIconGroup{}
	saveData.ID = req.ID
	saveData.CardStyle = base.AnyToJsonStr(req.CardStyle)
	saveData.Title = req.Title
	saveData.Icon = req.Icon
	saveData.OnlyName = req.OnlyName
	// saveData.Icon = "material-symbols:ad-group-outline"
	saveData.Description = req.Description
	saveData.Sort = req.Sort
	saveData.UserId = userInfo.ID
	saveData.PublicVisitHide = base.BoolToInt(req.PublicVisitHide)

	if v, err := biz.ItemGroup.Edit(userInfo.ID, saveData); err != nil {
		if errors.Is(err, biz.ErrItemOnlyNameExisted) {
			apiReturn.ErrorByCodeAndMsg(c, -2, err.Error())
			return
		}
		apiReturn.ErrorDatabase(c, err.Error())
		return
	} else {
		saveData = v
	}

	apiReturn.SuccessData(c, saveData)
}

func (a *ItemIconGroup) GetList(c *gin.Context) {

	userInfo, _ := base.GetCurrentUserInfo(c)
	currentMode := base.GetCurrentVisitMode(c)
	groups := []models.ItemIconGroup{}

	err := global.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Order("sort ,created_at").Where("user_id=?", userInfo.ID).Find(&groups).Error; err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
			return err
		}

		// 判断分组是否为空，为空将自动创建默认分组
		if len(groups) == 0 {
			defaultGroup := models.ItemIconGroup{
				Title:  "APP",
				UserId: userInfo.ID,
				Icon:   "material-symbols:ad-group-outline",
			}
			if err := tx.Create(&defaultGroup).Error; err != nil {
				apiReturn.ErrorDatabase(c, err.Error())
				return err
			}

			// 并将当前账号下所有无分组的图标更新到当前组
			if err := tx.Model(&models.ItemIcon{}).Where("user_id=?", userInfo.ID).Update("item_icon_group_id", defaultGroup.ID).Error; err != nil {
				apiReturn.ErrorDatabase(c, err.Error())
				return err
			}

			groups = append(groups, defaultGroup)
		}

		// 返回 nil 提交事务
		return nil
	})

	resp := []panelApiStructs.ItemGroupInfo{}
	for _, v := range groups {
		cardStyle := panelApiStructs.ItemGroupInfoCardStyle{}

		// 过滤公开隐藏的
		if currentMode == base.VISIT_MODE_PUBLIC && base.IntToBool(v.PublicVisitHide) {
			continue
		}

		base.JsonStrToAny(v.CardStyle, &cardStyle)
		resp = append(resp, panelApiStructs.ItemGroupInfo{
			ID:              v.ID,
			Icon:            v.Icon,
			Title:           v.Title,
			OnlyName:        v.OnlyName,
			Description:     v.Description,
			Sort:            v.Sort,
			UserId:          v.UserId,
			User:            v.User,
			CardStyle:       cardStyle,
			PublicVisitHide: base.IntToBool(v.PublicVisitHide),
		})
	}

	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	} else {
		apiReturn.SuccessListData(c, resp, 0)
	}
}

func (a *ItemIconGroup) Deletes(c *gin.Context) {
	req := commonApiStructs.RequestDeleteIds[uint]{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	userInfo, _ := base.GetCurrentUserInfo(c)

	var count int64
	if err := global.Db.Model(&models.ItemIconGroup{}).Where(" user_id=?", userInfo.ID).Count(&count).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	} else {
		if math.Abs(float64(len(req.Ids))-float64(count)) < 1 {
			apiReturn.ErrorCode(c, 1201, "At least one must be retained", nil)
			return
		}

	}

	txErr := global.Db.Transaction(func(tx *gorm.DB) error {
		mitemIcon := models.ItemIcon{}
		if err := tx.Delete(&models.ItemIconGroup{}, "id in ? AND user_id=?", req.Ids, userInfo.ID).Error; err != nil {
			return err
		}

		if err := mitemIcon.DeleteByItemIconGroupIds(tx, userInfo.ID, req.Ids); err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		apiReturn.ErrorDatabase(c, txErr.Error())
		return
	}

	apiReturn.Success(c)
}

// 保存排序
func (a *ItemIconGroup) SaveSort(c *gin.Context) {
	req := commonApiStructs.SortRequest{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)

	transactionErr := global.Db.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		for _, v := range req.SortItems {
			if err := tx.Model(&models.ItemIconGroup{}).Where("user_id=? AND id=?", userInfo.ID, v.Id).Update("sort", v.Sort).Error; err != nil {
				// 返回任何错误都会回滚事务
				return err
			}
		}

		// 返回 nil 提交事务
		return nil
	})

	if transactionErr != nil {
		apiReturn.ErrorDatabase(c, transactionErr.Error())
		return
	}

	apiReturn.Success(c)
}
