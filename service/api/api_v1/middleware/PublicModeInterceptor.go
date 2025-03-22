package middleware

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/lib/cmn/systemSetting"
	"sun-panel/models"

	"sync"

	"github.com/gin-gonic/gin"
)

var dbMutex sync.Mutex

// 公开访问模式（访客模式）
func PublicModeInterceptor(c *gin.Context) {
	// 获得token
	cToken := c.GetHeader("token")
	if cToken == "" {
		// 处理没有token的情况
		handleNoToken(c)
		return
	}

	var token string
	var ok bool

	if token, ok = global.CUserToken.Get(cToken); !ok || token == "" {
		// 处理token无效或为空的情况
		handleInvalidToken(c, cToken)
		return
	}

	// 尝试从缓存中获取用户信息
	if userInfo, success := global.UserToken.Get(token); success {
		global.Logger.Debug("缓存的用户TOKEN-1:", token)
		c.Set("userInfo", userInfo)
		return
	}

	// 如果缓存中没有，尝试获取锁，以确保只有一个请求去数据库查询
	dbMutex.Lock()
	defer dbMutex.Unlock()

	// 再次检查缓存，以防在等待锁的时候已经被其他请求填充了缓存
	if userInfo, success := global.UserToken.Get(token); success {
		global.Logger.Debug("缓存的用户TOKEN:", token)
		c.Set("userInfo", userInfo)
		return
	}

	// 缓存中没有，执行数据库查询
	mUser := models.User{}
	if info, err := mUser.GetUserInfoByToken(token); err == nil && info.Token != "" && info.ID != 0 {
		global.Logger.Debug("数据库查询用户:", info.ID)
		global.UserToken.SetDefault(info.Token, info)
		global.CUserToken.SetDefault(cToken, token)
		c.Set("userInfo", info)
	} else {
		global.Logger.Debug("数据库查询用户失败", token)
		apiReturn.ErrorByCode(c, apiReturn.ErrCodeTokenExpired)
		c.Abort()
	}
}

// 处理没有token的情况
// 获取公开账号的信息，如果存在则是开启公开模式
func handleNoToken(c *gin.Context) {
	var userId *uint
	if err := global.SystemSetting.GetValueByInterface(systemSetting.PANEL_PUBLIC_USER_ID, &userId); err == nil && userId != nil {
		userInfo := models.User{}
		// 从数据库获取公开账号的信息，这里应该加上缓存
		if err := global.Db.First(&userInfo, "id=?", userId).Error; err != nil {
			apiReturn.ErrorByCode(c, apiReturn.ErrCodeNotLoggedIn)
			c.Abort()
		} else {
			global.Logger.Debug("访客用户ID:", userInfo.ID)
			c.Set("userInfo", userInfo)
			c.Set(base.GIN_GET_VISIT_MODE, base.VISIT_MODE_PUBLIC)
		}
	} else {
		global.Logger.Debug("访客用户不存在")
		apiReturn.ErrorByCode(c, apiReturn.ErrCodeNotLoggedIn)
		c.Abort()
	}
}

// 处理token无效或为空的情况
func handleInvalidToken(c *gin.Context, cToken string) {
	global.Logger.Debug("token为空或者不OK")
	apiReturn.ErrorByCode(c, apiReturn.ErrCodeTokenExpired)
	c.Abort()
}
