package panel

import (
	"sun-panel/api/openAPI_v1/common/apiReturn"
	"sun-panel/lib/cmn"

	"github.com/gin-gonic/gin"
)

type Version struct {
}

func (a *Version) GetInfo(c *gin.Context) {
	info := cmn.GetSysVersionInfo()
	apiReturn.SuccessData(c, gin.H{
		"version":     info.Version,
		"versionCode": info.Version_code, // 辅助代号
	})
}
