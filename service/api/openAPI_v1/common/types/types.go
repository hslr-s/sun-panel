package types

type IconItem struct {
	// Type int `json:"type"` // 图标类型
	IconUrl     string `json:"iconUrl"`
	Title       string `json:"title"`
	OnlyName    string `json:"onlyName"`
	Url         string `json:"url"`
	LanUrl      string `json:"lanUrl"`
	Description string `json:"description"`
	// OpenMethod  int
	// Sort        int
	// ItemIconGroupId int 不可选择
}

type GetIconItemInfoReq struct {
	OnlyName string `json:"onlyName"`
}
