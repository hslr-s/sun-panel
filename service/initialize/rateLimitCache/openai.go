package rateLimitCache

import (
	"sun-panel/global"
	"sun-panel/lib/cache"
	"time"
)

// 速率限制分钟级别
func InitMinute() cache.Cacher[int] {
	return global.NewCache[int](1*time.Minute, 1*time.Hour, "RateLimitCacheMinute")
}

// 速率限制小时级别
func InitHour() cache.Cacher[int] {
	return global.NewCache[int](1*time.Hour, 2*time.Hour, "RateLimitCacheHour")
}
