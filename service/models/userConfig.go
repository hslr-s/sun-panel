package models

type UserConfig struct {
	UserId uint `gorm:"index" json:"userId"`

	// 纯前端数据，面板样式数据
	PanelJson string                 `json:"-"`
	Panel     map[string]interface{} `gorm:"-" json:"panel"`

	// 搜索引擎
	SearchEngineJson string                 `json:"-"`
	SearchEngine     map[string]interface{} `gorm:"-" json:"searchEngine"`
}
