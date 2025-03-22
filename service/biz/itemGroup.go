package biz

import (
	"sun-panel/global"
	"sun-panel/models"
)

type ItemGroupType struct {
}

func (p *ItemGroupType) Edit(userId uint, itemGroup models.ItemIconGroup) (models.ItemIconGroup, error) {

	itemGroup.UserId = userId

	if itemGroup.ID != 0 {

		if v, err := p.FindByOnlyName(userId, itemGroup.OnlyName); err == nil && v.ID != itemGroup.ID && itemGroup.OnlyName != "" {
			return itemGroup, ErrItemOnlyNameExisted
		}

		// 修改
		updateFields := []string{"Icon", "Title", "CardStyle", "UserId", "PublicVisitHide", "OnlyName"}

		global.Db.Model(&models.ItemIconGroup{}).
			Select(updateFields).
			Where("id=?", itemGroup.ID).Updates(&itemGroup)
	} else {
		if _, err := p.FindByOnlyName(userId, itemGroup.OnlyName); err == nil && itemGroup.OnlyName != "" {
			return itemGroup, ErrItemOnlyNameExisted
		}
		itemGroup.Sort = 9999
		// 创建
		global.Db.Create(&itemGroup)
	}

	return itemGroup, nil

}

func (p *ItemGroupType) FindByOnlyName(userId uint, onlyName string) (res models.ItemIconGroup, gormErr error) {
	gormErr = global.Db.First(&res, "only_name=? AND user_id=?", onlyName, userId).Error
	return
}

func (p *ItemGroupType) FindById(userId uint, id uint) (res models.ItemIconGroup, gormErr error) {
	gormErr = global.Db.First(&res, "id=? AND user_id=?", id, userId).Error
	return
}
