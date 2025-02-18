package global

import (
	"sun-panel/lib/cache"
	"time"
)

// 创建一个缓存区
// | defaultExpiration:默认过期时长
// | cleanupInterval:清理过期的key间隔 0.不清理
// | name:缓存名称
func NewCache[T any](defaultExpiration time.Duration, cleanupInterval time.Duration, name string) cache.Cacher[T] {
	return cache.NewGoCache[T](defaultExpiration, cleanupInterval)
}
