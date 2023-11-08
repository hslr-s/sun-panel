package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

// 系统设置
type SystemSetting struct {
	ID          uint   `gorm:"primaryKey"`
	ConfigName  string `gorm:"type:varchar(50)"`
	ConfigValue string `gorm:"type:text"`
}

func (m *SystemSetting) Get(configName string) (result string, err error) {
	var res SystemSetting
	if err := Db.Model(m).Select("ConfigValue").First(&res, "config_name=?", configName).Error; err != nil {
		return result, err
	}
	result = res.ConfigValue
	return result, nil
}

func (m *SystemSetting) GetValueByInterface(configName string, structValue interface{}) error {
	result, err := m.Get(configName)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(result), structValue)
	if err != nil {
		return err
	}
	return nil
}

// 暂时仅支持结构体和字符串
func (m *SystemSetting) Set(configName string, configValue interface{}) error {
	findRes := SystemSetting{}
	db := Db.Model(m).First(&findRes, "config_name=?", configName)
	if err := db.Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	value := ""
	if s, ok := configValue.(string); !ok {
		if jsonStr, err := json.Marshal(configValue); err != nil {
			value = ""
		} else {
			value = string(jsonStr)
		}
	} else {
		value = s
	}

	if db.RowsAffected == 0 {
		// 添加
		if err := Db.Model(m).Create(&SystemSetting{ConfigName: configName, ConfigValue: value}).Error; err != nil {
			return err
		}

	} else {
		// 修改
		if err := Db.Model(m).Where("id=?", findRes.ID).Update("config_value", value).Error; err != nil {
			return err
		}
	}
	return nil
}
