package panel

import (
	"errors"
	"strings"
	"sun-panel/api/openAPI_v1/common/apiReturn"
	"sun-panel/api/openAPI_v1/common/base"
	"sun-panel/api/openAPI_v1/common/types"
	"sun-panel/biz"
	"sun-panel/models"
	"sun-panel/models/datatype"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type Item struct {
}

func (a *Item) Create(c *gin.Context) {

	req := types.IconItem{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 获取用户的ID
	userInfo, _ := base.GetCurrentUserInfo(c)
	userId := userInfo.ID

	item := models.ItemIcon{
		Icon: datatype.ItemIconIconInfo{
			ItemType: datatype.ItemIconItemTypeImageUrl,
			Src:      req.IconUrl,
		},
		OnlyName:    strings.TrimSpace(req.OnlyName),
		Title:       strings.TrimSpace(req.Title),
		Url:         req.Url,
		LanUrl:      req.LanUrl,
		Description: req.Description,
		OpenMethod:  2,
	}

	if _, err := biz.Item.Edit(userId, item); err != nil {
		if errors.Is(err, biz.ErrItemOnlyNameExisted) {
			apiReturn.ErrorByCode(c, apiReturn.ErrCodeOnlyNameExisted)
			return
		}
		apiReturn.ErrorWithNum(c, err)
		return
	}

	apiReturn.Success(c)
}

func (a *Item) Update(c *gin.Context) {

	req := types.IconItem{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 获取用户的ID
	userInfo, _ := base.GetCurrentUserInfo(c)
	userId := userInfo.ID

	// 查询有没有该项目
	itemInfo, gormErr := biz.Item.FindByOnlyName(userId, req.OnlyName)

	if gormErr != nil {
		if errors.Is(gormErr, gorm.ErrRecordNotFound) {
			apiReturn.ErrorByCode(c, apiReturn.ErrCodeNoRecord)
			return
		}
		apiReturn.ErrorWithNum(c, gormErr)
		return
	}

	item := models.ItemIcon{
		Icon: datatype.ItemIconIconInfo{
			ItemType: datatype.ItemIconItemTypeImageUrl,
			Src:      req.IconUrl,
		},
		OnlyName:    strings.TrimSpace(req.OnlyName),
		Title:       strings.TrimSpace(req.Title),
		Url:         req.Url,
		LanUrl:      req.LanUrl,
		Description: req.Description,
		OpenMethod:  2,
	}

	item.ID = itemInfo.ID

	if _, err := biz.Item.Edit(userId, item); err != nil {
		if errors.Is(err, biz.ErrItemOnlyNameExisted) {
			apiReturn.ErrorByCode(c, apiReturn.ErrCodeOnlyNameExisted)
			return
		}
		apiReturn.ErrorWithNum(c, err)
		return
	}

	apiReturn.Success(c)
}

func (a *Item) GetInfoByOnlyName(c *gin.Context) {

	req := types.GetIconItemInfoReq{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 获取用户的ID
	userInfo, _ := base.GetCurrentUserInfo(c)
	userId := userInfo.ID

	req.OnlyName = strings.TrimSpace(req.OnlyName)

	info, gormErr := biz.Item.FindByOnlyName(userId, req.OnlyName)

	if gormErr != nil {
		if errors.Is(gormErr, gorm.ErrRecordNotFound) {
			apiReturn.ErrorByCode(c, apiReturn.ErrCodeNoRecord)
			return
		}
		apiReturn.ErrorWithNum(c, gormErr)
		return
	}

	resp := types.IconItem{
		OnlyName:    info.OnlyName,
		Title:       info.Title,
		Url:         info.Url,
		LanUrl:      info.LanUrl,
		Description: info.Description,
		IconUrl:     info.Icon.Src,
	}

	apiReturn.SuccessData(c, resp)
}
