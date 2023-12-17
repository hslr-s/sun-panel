package models

import (
	"sun-panel/models/datatype"

	"gorm.io/gorm"
)

type ItemIcon struct {
	BaseModel
	IconJson        string                    `gorm:"type:varchar(1000)" json:"-"`
	Icon            datatype.ItemIconIconInfo `gorm:"-" json:"icon"`
	Title           string                    `gorm:"type:varchar(50)" json:"title"`
	Url             string                    `gorm:"type:varchar(1000)" json:"url"`
	LanUrl          string                    `gorm:"type:varchar(1000)" json:"lanUrl"`
	Description     string                    `gorm:"type:varchar(1000)" json:"description"`
	OpenMethod      int                       `gorm:"type:tinyint(1)" json:"openMethod"`
	Sort            int                       `gorm:"type:int(11)" json:"sort"`
	ItemIconGroupId int                       `json:"itemIconGroupId"`
	UserId          uint                      `json:"userId"`
	User            User                      `json:"user"`
}

func (m *ItemIcon) DeleteByItemIconGroupIds(db *gorm.DB, userId uint, itemIconGroupIds []uint) (err error) {
	err = db.Delete(&ItemIcon{}, "item_icon_group_id in ? AND user_id=?", itemIconGroupIds, userId).Error
	return
}

func (m *ItemIcon) DeleteByUserId(db *gorm.DB, userId uint) (err error) {
	return db.Delete(&ItemIcon{}, "user_id=?", userId).Error
}
