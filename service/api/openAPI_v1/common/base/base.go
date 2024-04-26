package base

import (
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
)

func GetCurrentUserInfo(c *gin.Context) (userInfo models.User, exist bool) {

	if value, exist := c.Get("userInfo"); exist {
		if v, ok := value.(models.User); ok {
			return v, exist
		}
	}
	global.Logger.Debugln("GetCurrentUserInfo", userInfo.ID)
	return
}
