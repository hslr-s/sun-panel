package panelApiStructs

type GlobalSettingSite struct {
	FaviconUrl         string `json:"favicoUrl"`
	Title              string `json:"title"`
	LoginBackgroundUrl string `json:"loginBackgroundUrl"`
}

type GlobalSettingCustomJsAndCss struct {
	JsContent  string `json:"jsContent"`
	CssContent string `json:"cssContent"`
}
