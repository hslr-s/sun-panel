package api_v1

import (
	"sun-panel/api/api_v1/openness"
	"sun-panel/api/api_v1/panel"
	"sun-panel/api/api_v1/system"
)

type ApiGroup struct {
	ApiSystem system.ApiSystem // 系统功能api
	ApiOpen   openness.ApiPpenness
	ApiPanel  panel.ApiPanel
}

var ApiGroupApp = new(ApiGroup)
