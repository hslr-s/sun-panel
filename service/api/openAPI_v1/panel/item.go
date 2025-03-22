package panel

import (
	"errors"
	"fmt"
	"strings"
	"sun-panel/api/openAPI_v1/common/apiReturn"
	"sun-panel/api/openAPI_v1/common/base"
	"sun-panel/api/openAPI_v1/common/types"
	"sun-panel/biz"
	"sun-panel/global"
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

	item := models.ItemIcon{}
	item.OnlyName = req.OnlyName
	item.OpenMethod = 2 // 打开方式

	// 更新字段
	if req.IconUrl != nil {
		item.Icon.Src = *req.IconUrl
		item.Icon.ItemType = datatype.ItemIconItemTypeImageUrl
	}
	if req.Title != nil {
		item.Title = strings.TrimSpace(*req.Title)
	}
	if req.Url != nil {
		item.Url = *req.Url
	}
	if req.LanUrl != nil {
		item.LanUrl = *req.LanUrl
	}
	if req.Description != nil {
		item.Description = *req.Description
	}
	if req.ItemGroupID != nil || req.ItemGroupOnlyName != nil {
		groupInfo, groupOk := a.handleGetItemGroupInfoByIDOrOnlyName(c, userId, req.ItemGroupID, req.ItemGroupOnlyName)
		if !groupOk {
			return
		}
		item.ItemIconGroupId = int(groupInfo.ID)
	}
	// 保存图标并更新图标地址
	if req.IsSaveIcon != nil && *req.IsSaveIcon && req.IconUrl != nil {
		fileInfo, err := biz.File.DownloadFileAndSaveDb(global.Db, *req.IconUrl, userId)
		if err != nil {
			errmsg := "failed to save icon file :" + err.Error()
			apiReturn.Error(c, errmsg)
			// global.Logger.Error(errmsg)
			return
		}
		item.Icon.Src = global.UPLOAD_ROUTE + fileInfo.Src
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

	item := models.ItemIcon{}
	item = itemInfo
	// 更新字段
	if req.IconUrl != nil {
		item.Icon.Src = *req.IconUrl
		item.Icon.ItemType = datatype.ItemIconItemTypeImageUrl
	}
	if req.Title != nil {
		item.Title = strings.TrimSpace(*req.Title)
	}
	if req.Url != nil {
		item.Url = *req.Url
	}
	if req.LanUrl != nil {
		item.LanUrl = *req.LanUrl
	}
	if req.Description != nil {
		item.Description = *req.Description
	}
	if req.ItemGroupID != nil || req.ItemGroupOnlyName != nil {
		groupInfo, groupOk := a.handleGetItemGroupInfoByIDOrOnlyName(c, userId, req.ItemGroupID, req.ItemGroupOnlyName)
		if !groupOk {
			return
		}
		item.ItemIconGroupId = int(groupInfo.ID)
	}
	if item.OpenMethod == 0 {
		item.OpenMethod = 2 // 打开方式，新窗口打开
	}
	// 保存图标并更新图标地址
	if req.IsSaveIcon != nil && *req.IsSaveIcon && req.IconUrl != nil {
		fileInfo, err := biz.File.DownloadFileAndSaveDb(global.Db, *req.IconUrl, userId)
		if err != nil {
			errmsg := "failed to save icon file :" + err.Error()
			apiReturn.Error(c, errmsg)
			// global.Logger.Error(errmsg)
			return
		}
		item.Icon.Src = global.UPLOAD_ROUTE + fileInfo.Src
		item.Icon.ItemType = datatype.ItemIconItemTypeImageUrl
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

	itemGroupID := uint(info.ItemIconGroupId)
	groupInfo, _ := biz.ItemGroup.FindById(userId, itemGroupID)
	resp := types.IconItem{
		OnlyName:          info.OnlyName,
		Title:             &info.Title,
		Url:               &info.Url,
		LanUrl:            &info.LanUrl,
		Description:       &info.Description,
		IconUrl:           &info.Icon.Src,
		ItemGroupID:       &groupInfo.ID,
		ItemGroupOnlyName: &groupInfo.OnlyName,
	}

	apiReturn.SuccessData(c, resp)
}

// 根据ID或者onlyName查询项目分组信息，如果有错误直接接口返回，调用处无需再次返回接口错误
func (a *Item) handleGetItemGroupInfoByIDOrOnlyName(c *gin.Context, userId uint, groupID *uint, onlyName *string) (info models.ItemIconGroup, ok bool) {
	var (
		gormErr error
	)

	notGroupErrStr := "group information cannot be found"

	// 如果ID不为空,则根据ID查询
	if groupID != nil && *groupID != 0 {
		info, gormErr = biz.ItemGroup.FindById(userId, *groupID)
		if gormErr != nil {
			if errors.Is(gormErr, gorm.ErrRecordNotFound) {
				apiReturn.Error(c,
					fmt.Sprintf("%s [itemGroupID=%d]", notGroupErrStr, *groupID))
				return info, false
			}
			apiReturn.ErrorWithNum(c, gormErr)
			return info, false
		}
	} else if onlyName != nil && *onlyName != "" {
		info, gormErr = biz.ItemGroup.FindByOnlyName(userId, *onlyName)
		if gormErr != nil {
			if errors.Is(gormErr, gorm.ErrRecordNotFound) {
				apiReturn.Error(c,
					fmt.Sprintf("%s [itemGroupOnlyName=%s]", notGroupErrStr, *onlyName))
				return info, false
			}
			apiReturn.ErrorWithNum(c, gormErr)
			return info, false
		}
	} else {
		apiReturn.Error(c, notGroupErrStr)
		return info, false
	}

	return info, true
}
