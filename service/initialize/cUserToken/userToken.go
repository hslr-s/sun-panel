package cUserToken

import (
	"sun-panel/global"
	"sun-panel/lib/cache"

	"time"
)

func InitCUserToken(expiration time.Duration) cache.Cacher[string] {
	return global.NewCache[string](expiration, 48*time.Hour, "CUserToken")
}

// func InitVerifyCodeCachePool() {
// 	global.VerifyCodeCachePool = cache.NewGoCache(10*time.Minute, 60*time.Second)
// }
