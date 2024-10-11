package middleware

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
)

func LoginInterceptor(c *gin.Context) {

	// 继续执行后续的操作，再回来
	// c.Next()

	// 获得token
	cToken := c.GetHeader("token")

	// 没有token信息视为未登录
	if cToken == "" {
		apiReturn.ErrorByCode(c, apiReturn.ErrCodeNotLoggedIn)
		c.Abort() // 终止执行后续的操作，一般配合return使用
		return
	}

	token := ""
	{
		var ok bool
		token, ok = global.CUserToken.Get(cToken)
		// 可能已经安全退出或者很久没有使用已过期
		if !ok || token == "" {
			apiReturn.ErrorByCode(c, apiReturn.ErrCodeTokenExpired)
			c.Abort() // 终止执行后续的操作，一般配合return使用
			return
		}
	}

	// 直接返回缓存的用户信息
	if userInfo, success := global.UserToken.Get(token); success {
		c.Set("userInfo", userInfo)
		return
	}

	global.Logger.Debug("准备查询数据库的用户资料", token)

	mUser := models.User{}
	// 去库中查询是否存在该用户；否则返回错误
	if info, err := mUser.GetUserInfoByToken(token); err != nil || info.Token == "" || info.ID == 0 {
		apiReturn.ErrorByCode(c, apiReturn.ErrCodeTokenExpired)
		c.Abort()
		return
	} else {
		// 通过 设置当前用户信息 (这里可能存在了自动续期的功能)
		global.UserToken.SetDefault(info.Token, info)
		global.CUserToken.SetDefault(cToken, token)
		c.Set("userInfo", info)
	}

}

// 不验证缓存直接验证库省去没有缓存每次都要手动登录的问题
func LoginInterceptorDev(c *gin.Context) {

	// 获得token
	token := c.GetHeader("token")
	mUser := models.User{}

	// 去库中查询是否存在该用户；否则返回错误
	if info, err := mUser.GetUserInfoByToken(token); err != nil || info.ID == 0 {
		apiReturn.ErrorByCode(c, apiReturn.ErrCodeTokenExpired)
		c.Abort()
		return
	} else {
		// 通过
		// 设置当前用户信息
		c.Set("userInfo", info)
	}
}

// 调试模式登录-直接输入用户id模拟登录（不安全）
func LoginInterceptorDebug(c *gin.Context) {
	mUser := models.User{}
	// 去库中查询是否存在该用户；否则返回错误
	if info, err := mUser.GetUserInfoByUid(1); err != nil || info.ID == 0 {
		apiReturn.ErrorByCode(c, apiReturn.ErrCodeTokenExpired)
		c.Abort()
		return
	} else {
		// 通过
		// 设置当前用户信息
		c.Set("userInfo", info)
	}
}
