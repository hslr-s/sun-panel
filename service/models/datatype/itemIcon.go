package datatype

const (
	ItemIconItemTypeText       = 1
	ItemIconItemTypeImageUrl   = 2
	ItemIconItemTypeIconOnLine = 3
)

type ItemIconIconInfo struct {
	ItemType int    `json:"itemType"`
	Src      string `json:"src"`
	Text     string `json:"text"`

	// BgColor  string `json:"bgColor"`
	BackgroundColor string `json:"backgroundColor"`
}
