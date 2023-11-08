package jsonConfig

// 事件风格数据
type EventStyleDataModel struct {
	Title           string `json:"title"`           // 标题
	ClassName       string `json:"className"`       // 类名称
	TextColor       string `json:"textColor"`       // 字体颜色
	BackgroundColor string `json:"backgroundColor"` // 背景颜色
	BorderColor     string `json:"borderColor"`     // 边框颜色
}

type EventStyleModel struct {
	ConfigModel
	Data []EventStyleDataModel
}

func (e *EventStyleModel) GetImportData() error {
	return nil
}
