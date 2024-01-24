package cache

import (
	"time"
)

const (
	CACHE_DRIVE_REDIS  = "redis"
	CACHE_DRIVE_MEMORY = "memory"
)

// 缓存接口-支持Redis和内存使用
type Cacher[T any] interface {
	// 设置
	Set(k string, v T, d time.Duration)

	// 取值
	Get(k string) (T, bool)

	// 设置-过期时间采用默认值
	SetDefault(k string, v T)

	// 删除
	Delete(k string)

	// 只有在给定Key项尚未存在，或者现有项已过期时，才能将项添加到缓存中。否则返回错误。
	// Add(k string, v T, d time.Duration)
	// IncrementInt(k string, n int) (num int, err error)

	// 设置值，但不重置过期时间
	SetKeepExpiration(k string, v T)

	// 项目总数
	ItemCount() (int64, error)

	// 清空
	Flush()
}
