package panelApiStructs

import "sun-panel/biz"

type GlobalSettingSite struct {
	biz.GlobalSettingSiteSetting
}

type GlobalSettingCustomJsAndCss struct {
	JsContent  string `json:"jsContent"`
	CssContent string `json:"cssContent"`
}
