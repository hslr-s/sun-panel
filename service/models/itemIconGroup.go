package models

import (
	"gorm.io/gorm"
)

type ItemIconGroup struct {
	BaseModel
	Icon            string `json:"icon"`
	Title           string `gorm:"type:varchar(50)" json:"title"`
	OnlyName        string `gorm:"type:varchar(50);index" json:"onlyName"` // 唯一标识
	Description     string `gorm:"type:varchar(1000)" json:"description"`
	Sort            int    `gorm:"type:int(11)" json:"sort"`
	PublicVisitHide int    `gorm:"type:tinyint(1)" json:"publicVisitHide"` // 公开模式隐藏
	UserId          uint   `json:"userId"`
	User            User   `json:"user"`
	CardStyle       string `json:"cardStyle"`
}

func (m *ItemIconGroup) DeleteByUserId(db *gorm.DB, userId uint) (err error) {
	err = db.Delete(&ItemIconGroup{}, "user_id = ?", userId).Error
	return
}
