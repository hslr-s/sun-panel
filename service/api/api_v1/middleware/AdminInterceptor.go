package middleware

import (
	// "calendar-note-gin/api/v1/common/apiReturn"
	// . "calendar-note-gin/api/v1/common/base"

	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"

	"github.com/gin-gonic/gin"
)

func AdminInterceptor(c *gin.Context) {
	currentUser, _ := base.GetCurrentUserInfo(c)
	global.Logger.Debugln("登录用户", currentUser)
	if currentUser.Role != 1 {
		apiReturn.ErrorNoAccess(c)
		c.Abort()
		return
	}
}
