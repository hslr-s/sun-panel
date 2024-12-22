package panel

import (
	"errors"
	"strings"
	"sun-panel/api/openAPI_v1/common/apiReturn"
	"sun-panel/api/openAPI_v1/common/base"
	"sun-panel/api/openAPI_v1/common/types"
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

// 项目组
type ItemGroup struct {
}

func (a *ItemGroup) Create(c *gin.Context) {

	req := types.ItemGroup{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 获取用户的ID
	userInfo, _ := base.GetCurrentUserInfo(c)
	userId := userInfo.ID

	itemGroup := models.ItemIconGroup{}
	itemGroup.OnlyName = req.OnlyName

	if req.Title != nil {
		itemGroup.Title = strings.TrimSpace(*req.Title)
	}

	if _, err := biz.ItemGroup.Edit(userId, itemGroup); err != nil {
		if errors.Is(err, biz.ErrItemOnlyNameExisted) {
			apiReturn.ErrorByCode(c, apiReturn.ErrCodeOnlyNameExisted)
			return
		}
		apiReturn.ErrorWithNum(c, err)
		return
	}

	apiReturn.Success(c)
}

// func (a *ItemGroup) GetInfoByOnlyName(c *gin.Context) {

// 	req := types.GetItemGroupInfoByOnleyNameReq{}

// 	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
// 		apiReturn.ErrorParamFomat(c, err.Error())
// 		return
// 	}

// 	// 获取用户的ID
// 	userInfo, _ := base.GetCurrentUserInfo(c)
// 	userId := userInfo.ID

// 	req.OnlyName = strings.TrimSpace(req.OnlyName)

// 	info, gormErr := biz.ItemGroup.FindByOnlyName(userId, req.OnlyName)

// 	if gormErr != nil {
// 		if errors.Is(gormErr, gorm.ErrRecordNotFound) {
// 			apiReturn.ErrorByCode(c, apiReturn.ErrCodeNoRecord)
// 			return
// 		}
// 		apiReturn.ErrorWithNum(c, gormErr)
// 		return
// 	}

// 	resp := types.ItemGroup{
// 		OnlyName:    info.OnlyName,
// 		Title:       &info.Title,
// 		ItemGroupID: &info.ID,
// 	}

// 	apiReturn.SuccessData(c, resp)
// }

func (a *ItemGroup) GetInfo(c *gin.Context) {

	req := types.GetItemGroupInfoReq{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 获取用户的ID
	userInfo, _ := base.GetCurrentUserInfo(c)
	userId := userInfo.ID

	var (
		info    models.ItemIconGroup
		gormErr error
	)

	// 如果ID不为空,则根据ID查询
	if req.ItemGroupID != nil && *req.ItemGroupID != 0 {
		info, gormErr = biz.ItemGroup.FindById(userId, *req.ItemGroupID)

		if gormErr != nil {
			if errors.Is(gormErr, gorm.ErrRecordNotFound) {
				apiReturn.ErrorByCode(c, apiReturn.ErrCodeNoRecord)
				return
			}
			apiReturn.ErrorWithNum(c, gormErr)
			return
		}
	} else if req.OnlyName != nil && *req.OnlyName != "" {
		info, gormErr = biz.ItemGroup.FindByOnlyName(userId, *req.OnlyName)
		if gormErr != nil {
			if errors.Is(gormErr, gorm.ErrRecordNotFound) {
				apiReturn.ErrorByCode(c, apiReturn.ErrCodeNoRecord)
				return
			}
			apiReturn.ErrorWithNum(c, gormErr)
			return
		}
	} else {
		apiReturn.ErrorByCode(c, apiReturn.ErrCodeNoRecord)
		return
	}

	resp := types.ItemGroup{
		OnlyName:    info.OnlyName,
		Title:       &info.Title,
		ItemGroupID: &info.ID,
	}

	apiReturn.SuccessData(c, resp)
}

// 查询所有的项目分组
func (a *ItemGroup) GetList(c *gin.Context) {
	// 获取用户的ID
	userInfo, _ := base.GetCurrentUserInfo(c)
	groups := []models.ItemIconGroup{}
	err := global.Db.Order("sort,created_at").
		Where("user_id=?", userInfo.ID).
		Find(&groups).Error

	if err != nil {
		apiReturn.ErrorWithNum(c, err)
		return
	}

	respList := []types.ItemGroup{}
	for _, v := range groups {
		title := v.Title
		itemGroupId := v.ID

		respList = append(respList, types.ItemGroup{
			OnlyName:    v.OnlyName,
			Title:       &title,
			ItemGroupID: &itemGroupId,
		})
	}

	apiReturn.SuccessListData(c, respList, int64(len(respList)))
}
