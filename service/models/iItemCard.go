package models

import (
	"gorm.io/gorm"
)

// 【暂未使用】从1.5.0+ 推荐模型
type ItemCard struct {
	BaseModel
	CardType        int    `gorm:"type:tinyint(1)" json:"cardType"` // 卡片类型：普通图标，docker图标，系统统计卡片（暂时不适配）...
	Title           string `gorm:"type:varchar(50)" json:"title"`
	OnlyName        string `gorm:"type:varchar(200)" json:"onlyName"`
	Description     string `gorm:"type:varchar(1000)" json:"description"`
	Sort            int    `gorm:"type:int(11)" json:"sort"`
	Background      string `gorm:"type:varchar(20)" json:"background"` // 背景色 1.5.0
	ItemIconGroupId int    `json:"itemIconGroupId"`
	ExpandParam     string `json:"expandParam"`
	UserId          uint   `json:"userId"`
	User            User   `json:"user"`
}

func (m *ItemCard) DeleteByItemCardGroupIds(db *gorm.DB, userId uint, itemIconGroupIds []uint) (err error) {
	err = db.Delete(&ItemCard{}, "item_icon_group_id in ? AND user_id=?", itemIconGroupIds, userId).Error
	return
}

func (m *ItemCard) DeleteByUserId(db *gorm.DB, userId uint) (err error) {
	return db.Delete(&ItemCard{}, "user_id=?", userId).Error
}
