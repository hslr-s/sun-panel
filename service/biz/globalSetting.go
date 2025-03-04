package biz

import (
	"errors"
	"sun-panel/global"
)

type GlobalSettingType struct {
}

type GlobalSettingSiteSetting struct {
	FaviconUrl             string `json:"favicoUrl"`
	Title                  string `json:"title"`
	LoginBackgroundUrl     string `json:"loginBackgroundUrl"`
	IsEnableLoginCaptcha   bool   `json:"isEnableLoginCaptcha"`
	IsDisableItemCardCache bool   `json:"isDisableItemCardCache"`
}

var ErrGlobalSettingTypeNoPROAuthorization = errors.New("no PRO authorization")

func (g *GlobalSettingType) SaveSiteStting(value GlobalSettingSiteSetting) error {
	return global.SystemSetting.Set("global_setting_site", value)
}

// value传入指针接收
func (g *GlobalSettingType) GetSiteStting(value interface{}) error {
	if ProAuth.ProIsExpired() {
		return ErrGlobalSettingTypeNoPROAuthorization
	}
	return global.SystemSetting.GetValueByInterface("global_setting_site", value)
}

// value传入指针接收
func (g *GlobalSettingType) GetSiteSttingNoCheckExpired(value interface{}) error {
	return global.SystemSetting.GetValueByInterface("global_setting_site", value)
}
