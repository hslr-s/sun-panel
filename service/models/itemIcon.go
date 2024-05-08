package models

import (
	"sun-panel/models/datatype"

	"gorm.io/gorm"
)

// 从1.5.0+ 此表别名为 card
type ItemIcon struct {
	BaseModel
	CardType        int                       `gorm:"type:tinyint(1)" json:"cardType"` // 1.5.0新增 // 卡片类型：普通图标，docker图标，系统统计卡片（暂时不适配）...
	IconJson        string                    `gorm:"type:varchar(1000)" json:"-"`     // 历史遗留问题，为适配1.4.0及以下版本暂时适配
	Icon            datatype.ItemIconIconInfo `gorm:"-" json:"icon"`
	Title           string                    `gorm:"type:varchar(50)" json:"title"` // 1.4.0新增
	OnlyName        string                    `gorm:"type:varchar(200)" json:"onlyName"`
	Url             string                    `gorm:"type:varchar(1000)" json:"url"`
	LanUrl          string                    `gorm:"type:varchar(1000)" json:"lanUrl"`
	Description     string                    `gorm:"type:varchar(1000)" json:"description"`
	OpenMethod      int                       `gorm:"type:tinyint(1)" json:"openMethod"`
	Sort            int                       `gorm:"type:int(11)" json:"sort"`
	BackgroundColor string                    `gorm:"type:varchar(20)" json:"backgroundColor"` // 背景色 1.5.0新增
	ItemIconGroupId int                       `json:"itemIconGroupId"`
	ExpandParam     string                    `json:"expandParam"` // 1.5.0新增
	UserId          uint                      `json:"userId"`
	User            User                      `json:"user"`
	// Usersss         string                    `json:"usersss"`
}

func (m *ItemIcon) DeleteByItemIconGroupIds(db *gorm.DB, userId uint, itemIconGroupIds []uint) (err error) {
	err = db.Delete(&ItemIcon{}, "item_icon_group_id in ? AND user_id=?", itemIconGroupIds, userId).Error
	return
}

func (m *ItemIcon) DeleteByUserId(db *gorm.DB, userId uint) (err error) {
	return db.Delete(&ItemIcon{}, "user_id=?", userId).Error
}
