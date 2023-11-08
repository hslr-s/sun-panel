package iniConfig

import (
	"gopkg.in/ini.v1"
)

type IniConfig struct {
	Err      error
	Config   *ini.File
	Default  map[string]map[string]string // 默认配置
	FileName string
}

// 获取配置
func (t *IniConfig) GetValue(section string, name string) *ini.Key {
	return t.Config.Section(section).Key(name)
}

// 设置配置
func (t *IniConfig) SetValue(section string, name string, value string) error {
	t.Config.Section(section).Key(name).SetValue(value)
	return t.Config.SaveTo(t.FileName)
}

// 获取配置
func (t *IniConfig) GetValueString(section string, name string) string {
	return t.Config.Section(section).Key(name).String()
}

// 获取字符串配置，如果没有会查找默认值
func (t *IniConfig) GetValueStringOrDefault(section string, name string) string {
	value := t.GetValueString(section, name)
	if value == "" && t.Default[section] != nil && t.Default[section][name] != "" {
		return t.Default[section][name]
	} else {
		return value
	}
}

// 获取配置
func (t *IniConfig) GetValueInt(section string, name string) int {
	return t.Config.Section(section).Key(name).MustInt()
}

// 获取组配置
func (t *IniConfig) GetSection(section string, result interface{}) error {
	if group, err := t.Config.GetSection(section); err != nil {
		return err
	} else {
		if err := group.MapTo(result); err != nil {
			return err
		} else {
			return nil
		}
	}
}

// 删除组
func (t *IniConfig) DeleteSection(section string) {
	t.Config.DeleteSection(section)
	t.Config.SaveTo(t.FileName)
}

// 创建一个配置对象
func NewIniConfig(filename string) *IniConfig {
	config, err := ini.Load(filename)

	return &IniConfig{
		Err:      err,
		Config:   config,
		FileName: filename,
	}
}
