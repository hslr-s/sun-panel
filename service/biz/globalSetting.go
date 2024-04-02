package biz

import "sun-panel/global"

type GlobalSettingType struct {
}

type GlobalSettingSiteSetting struct {
	FaviconUrl         string `json:"favicoUrl"`
	Title              string `json:"title"`
	LoginBackgroundUrl string `json:"loginBackgroundUrl"`
}

func (g *GlobalSettingType) SaveSiteStting(value GlobalSettingSiteSetting) error {
	return global.SystemSetting.Set("global_setting_site", value)
}

// value传入指针接收
func (g *GlobalSettingType) GetSiteStting(value interface{}) error {
	return global.SystemSetting.GetValueByInterface("global_setting_site", value)
}
