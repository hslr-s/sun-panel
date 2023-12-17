package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

type ModuleConfig struct {
	BaseModel
	UserId    uint                   `gorm:"index" json:"userId"`
	Name      string                 `gorm:"type:varchar(255)" json:"name"`
	ValueJson string                 `gorm:"type:text" json:"-"`
	Value     map[string]interface{} `gorm:"-" json:"value"`
}

func (m *ModuleConfig) GetConfigByUserIdAndName(db *gorm.DB, userId uint, name string) (map[string]interface{}, error) {
	cfg := ModuleConfig{}
	if err := db.First(&cfg, "user_id=? AND name=?", userId, name).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}

	// 处理字段
	if err := json.Unmarshal([]byte(cfg.ValueJson), &cfg.Value); err != nil {
		cfg.Value = nil
	}
	return cfg.Value, nil
}

func (m *ModuleConfig) Save(db *gorm.DB) error {

	// 处理字段
	if jb, err := json.Marshal(m.Value); err != nil {
		m.ValueJson = "{}"
	} else {
		m.ValueJson = string(jb)
	}

	// 保存操作
	if err := db.First(&ModuleConfig{}, "user_id=? AND name=?", m.UserId, m.Name).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 新增
			if err := db.Create(&m).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		// 修改
		if err := db.Select("Name", "UserId", "ValueJson").Where("user_id=? AND name=?", m.UserId, m.Name).Updates(&m).Error; err != nil {
			return err
		}
	}

	return nil
}
