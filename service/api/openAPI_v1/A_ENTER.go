package openAPI_v1

import "sun-panel/api/openAPI_v1/panel"

type OpenApiGroup struct {
	PanelApiGroup panel.PanelApiGroup
}

var OpenApiGroupApp = new(OpenApiGroup)
