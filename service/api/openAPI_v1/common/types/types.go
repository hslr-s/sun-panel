package types

type ItemGroup struct {
	OnlyName    string  `json:"onlyName"`
	ItemGroupID *uint   `json:"itemGroupID"`
	Title       *string `json:"title"`
}

type IconItem struct {
	// Type int `json:"type"` // 图标类型
	IconUrl           *string `json:"iconUrl"`
	Title             *string `json:"title"`
	OnlyName          string  `json:"onlyName"`
	Url               *string `json:"url"`
	LanUrl            *string `json:"lanUrl"`
	Description       *string `json:"description"`
	ItemGroupID       *uint   `json:"itemGroupID"`       // 组ID v1.7.* 新增
	ItemGroupOnlyName *string `json:"itemGroupOnlyName"` // 组名字 v1.7.* 新增
	IsSaveIcon        *bool   `json:"isSaveIcon"`        // 是否保存图标,保存图标将进行下载 v1.7.* 新增
	// OpenMethod  int
	// Sort        int
	// ItemIconGroupId int 不可选择
}

type GetIconItemInfoReq struct {
	OnlyName string `json:"onlyName"`
}

// type GetItemGroupInfoByItemGroupIDReq struct {
// 	ItemGroupID uint `json:"itemGroupID"`
// }

// type GetItemGroupInfoByOnleyNameReq struct {
// 	OnlyName string `json:"onlyName"`
// }

type GetItemGroupInfoReq struct {
	OnlyName    *string `json:"onlyName"`
	ItemGroupID *uint   `json:"itemGroupID"`
}
