package panelApiStructs

import (
	"sun-panel/models"
)

type ItemGroupInfoCardStyle struct {
	Style                   int    `json:"style"`
	TextColor               string `json:"textColor"`
	TextInfoHideDescription bool   `json:"textInfoHideDescription"`
	TextIconHideTitle       bool   `json:"textIconHideTitle"`
}

type ItemGroupInfo struct {
	ID              uint                   `json:"id"`
	Icon            string                 `json:"icon"`
	Title           string                 `json:"title"`
	OnlyName        string                 `json:"onlyName"`
	Description     string                 `json:"description"`
	Sort            int                    `json:"sort"`
	UserId          uint                   `json:"userId"`
	PublicVisitHide bool                   `json:"publicVisitHide"`
	User            models.User            `json:"user"`
	CardStyle       ItemGroupInfoCardStyle `json:"cardStyle"`
}

type ItemGroupEditReq struct {
	ItemGroupInfo
}
