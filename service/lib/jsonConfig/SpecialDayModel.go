package jsonConfig

type SpecialDayModel struct {
	ConfigModel
	Data SpecialDayDataModel
}

type SpecialDayDataModel struct {
	OnlyId     string `json:"onlyId"`     // 唯一ID
	Name       string `json:"name"`       // 名称
	UpdateTime string `json:"updateTime"` // 更新时间
	// Year       int                                `json:"year"`       // 年份
	StartDate string                             `json:"startDate"` // 开始日期
	EndDate   string                             `json:"endDate"`   // 结束日期
	Days      map[string]SpecialDayDataDaysModel `json:"days"`      // 天数据
}

type SpecialDayDataDaysModel struct {
	Holiday bool   `json:"holiday"` // 唯一ID
	Name    string `json:"name"`    // 名称
}
