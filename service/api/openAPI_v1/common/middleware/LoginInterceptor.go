package middleware

import (
	"errors"
	"sun-panel/api/openAPI_v1/common/apiReturn"
	"sun-panel/biz"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LoginInterceptor(c *gin.Context) {

	// 获得token
	token := c.GetHeader("token")
	// global.Logger.Debug("token:", token)

	// 没有token信息视为未登录
	if token == "" {
		apiReturn.ErrorByCode(c, apiReturn.ErrCodeTokenExpires)
		c.Abort() // 终止执行后续的操作，一般配合return使用
		return
	}

	if OpenApiToken, gormErr := biz.OpenApi.GetUserInfo(token); gormErr != nil {
		if errors.Is(gormErr, gorm.ErrRecordNotFound) {
			apiReturn.ErrorByCode(c, apiReturn.ErrCodeTokenExpires)
			c.Abort() // 终止执行后续的操作，一般配合return使用
			return
		}
		apiReturn.ErrorWithNum(c, gormErr)
		c.Abort() // 终止执行后续的操作，一般配合return使用
		return
	} else {
		c.Set("userInfo", OpenApiToken.User)
		// global.Logger.Debugln("userInfo", OpenApiToken.User.ID)
	}
}
