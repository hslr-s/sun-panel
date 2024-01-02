package system

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/global"

	"github.com/gin-gonic/gin"
)

type MonitorApi struct{}

func (a *MonitorApi) GetAll(c *gin.Context) {
	if value, ok := global.SystemMonitor.Get("value"); ok {
		apiReturn.SuccessData(c, value)
		return
	}
	apiReturn.Error(c, "failed")
}
