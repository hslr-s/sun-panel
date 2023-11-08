package userToken

import (
	"sun-panel/global"
	"sun-panel/lib/cache"
	"sun-panel/models"

	"time"
)

func InitUserToken() cache.Cacher[models.User] {
	return global.NewCache[models.User](1*time.Minute, 1*time.Hour, "UserToken")
}

// func InitVerifyCodeCachePool() {
// 	global.VerifyCodeCachePool = cache.NewGoCache(10*time.Minute, 60*time.Second)
// }
