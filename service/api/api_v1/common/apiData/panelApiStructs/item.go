package panelApiStructs

import (
	"sun-panel/api/api_v1/common/apiData/commonApiStructs"
	"sun-panel/models"
	"sun-panel/models/datatype"
)

type ItemInfo struct {
	models.BaseModel
	CardType        int                       `json:"cardType"` // 卡片类型：普通图标，docker图标，系统统计卡片（暂时不适配）...
	IconJson        string                    `json:"-"`        // 历史遗留问题，为适配1.4.0及以下版本暂时适配
	Icon            datatype.ItemIconIconInfo `json:"icon"`
	Title           string                    `json:"title"`
	OnlyName        string                    `json:"onlyName"`
	Url             string                    `json:"url"`
	LanUrl          string                    `json:"lanUrl"`
	Description     string                    `json:"description"`
	OpenMethod      int                       `json:"openMethod"`
	Sort            int                       `json:"sort"`
	BackgroundColor string                    `json:"backgroundColor"` // 背景色 1.5.0
	ItemIconGroupId int                       `json:"itemIconGroupId"`
	ExpandParam     map[string]any            `json:"expandParam"`
	UserId          uint                      `json:"userId"`
}

type ItemIconEditRequest struct {
	models.ItemIcon
	IconJson string
}

type ItemIconSaveSortRequest struct {
	SortItems       []commonApiStructs.SortRequestItem `json:"sortItems"`
	ItemIconGroupId uint                               `json:"itemIconGroupId"`
}

type ItemIconGetSiteFaviconReq struct {
	Url string `json:"url"`
}

type ItemIconGetSiteFaviconResp struct {
	IconUrl string `json:"iconUrl"`
}
