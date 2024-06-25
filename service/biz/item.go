package biz

import (
	"encoding/json"
	"errors"
	"sun-panel/global"
	"sun-panel/models"
)

type ItemType struct {
}

var (
	ErrItemOnlyNameExisted = errors.New("only name existed")
)

func (p *ItemType) Edit(userId uint, itemIcon models.ItemIcon) (models.ItemIcon, error) {

	// 没有分组自行选择一个排在第一位的分组
	if itemIcon.ItemIconGroupId == 0 {
		defaultGroup := models.ItemIconGroup{}
		global.Db.Order("sort ,created_at").Where("user_id=?", userId).First(&defaultGroup)
		itemIcon.ItemIconGroupId = int(defaultGroup.ID)
	}

	itemIcon.UserId = userId

	// json转字符串
	if j, err := json.Marshal(itemIcon.Icon); err == nil {
		itemIcon.IconJson = string(j)
	}

	if itemIcon.ID != 0 {

		if v, err := p.FindByOnlyName(userId, itemIcon.OnlyName); err == nil && v.ID != itemIcon.ID && itemIcon.OnlyName != "" {
			return itemIcon, ErrItemOnlyNameExisted
		}

		// 修改
		updateField := []string{
			"IconJson",
			"Icon",
			"Title",
			"Url",
			"LanUrl",
			"Description",
			"OpenMethod",
			"GroupId",
			"UserId",
			"ItemIconGroupId",
			"OnlyName",
			"CardType",
			"BackgroundColor",
			"ExpandParam",
		}
		// 修改不更新sort
		// if itemIcon.Sort != 0 {
		// 	updateField = append(updateField, "Sort")
		// }
		global.Db.Model(&models.ItemIcon{}).
			Select(updateField).
			Where("id=?", itemIcon.ID).Updates(&itemIcon)
	} else {
		if _, err := p.FindByOnlyName(userId, itemIcon.OnlyName); err == nil && itemIcon.OnlyName != "" {
			return itemIcon, ErrItemOnlyNameExisted
		}
		itemIcon.Sort = 9999
		// 创建
		global.Db.Create(&itemIcon)
	}

	return itemIcon, nil

}

func (p *ItemType) FindByOnlyName(userId uint, onlyName string) (res models.ItemIcon, gormErr error) {
	gormErr = global.Db.First(&res, "only_name=? AND user_id=?", onlyName, userId).Error
	return
}
