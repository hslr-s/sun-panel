package cache

import (
	"context"
	"encoding/json"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type RedisCacheStruct[T any] struct {
	Redis             *redis.Client
	Ctx               context.Context
	HashKey           string
	Result            T
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration
}

type RedisValue[T any] struct {
	ExpirationTimeStamp int64
	IsExpiration        bool // 是否有过期时间 false // 不过期
	Value               T
}

// cache.New(5*time.Minute, 60*time.Second)
func NewRedisCache[T any](redisDb *redis.Client, hashKey string, defaultExpiration time.Duration, cleanupInterval time.Duration) *RedisCacheStruct[T] {
	obj := RedisCacheStruct[T]{
		Redis:             redisDb,
		Ctx:               context.Background(),
		HashKey:           hashKey,
		DefaultExpiration: defaultExpiration,
		CleanupInterval:   cleanupInterval,
	}

	// 创建定时器判断是否过期
	if obj.CleanupInterval.Seconds() > 0 {
		go obj.expirationVerification()
	}

	return &obj
}

func (r *RedisCacheStruct[T]) Set(k string, v T, d time.Duration) {
	valueEncode := ""
	value := RedisValue[T]{}

	// 设置过期时间
	if d.Seconds() > 0 {
		value.IsExpiration = true
		value.ExpirationTimeStamp = time.Now().Add(d).Unix()
	} else {
		value.IsExpiration = false // 不过期
	}

	value.Value = v
	if j, e := json.Marshal(value); e == nil {
		valueEncode = string(j)
	}

	r.Redis.HSet(r.Ctx, r.HashKey, k, valueEncode)
	// second := d.Seconds()
	// if second > 0 {
	// 	// 设置过期时间
	// 	err := r.Redis.Do(r.Ctx, "SETEX", r.HashKey+k, second, valueEncode).Err()
	// 	fmt.Println("设置结果", err)
	// } else {
	// 	r.Redis.HSet(r.Ctx, r.HashKey, k, valueEncode)
	// }

}

func (r *RedisCacheStruct[T]) Get(k string) (T, bool) {
	var valueEncode []byte
	value := RedisValue[T]{}
	cmd := r.Redis.HGet(r.Ctx, r.HashKey, k)
	if err := cmd.Scan(&valueEncode); err != nil {
		// log.Println(err)
		return r.Result, false
	}

	if err := json.Unmarshal(valueEncode, &value); err != nil {
		// log.Println(err)
		return r.Result, false
	}

	// 已过期，清理掉key
	if value.IsExpiration && time.Now().Unix() > value.ExpirationTimeStamp {
		r.Delete(k)
		return r.Result, false
	}

	return value.Value, true
}

// 设置cache 无时间参数
func (r *RedisCacheStruct[T]) SetDefault(k string, v T) {
	r.Set(k, v, r.DefaultExpiration)
}

// 设置并保持原始的过期时间
func (r *RedisCacheStruct[T]) SetKeepExpiration(k string, v T) {
	var valueEncode []byte
	value := RedisValue[T]{}
	cmd := r.Redis.HGet(r.Ctx, r.HashKey, k)
	if err := cmd.Scan(&valueEncode); err != nil {
		// fmt.Println("使用默认的过期时间")
		r.SetDefault(k, v)
		return
	}

	if err := json.Unmarshal(valueEncode, &value); err != nil {
		// fmt.Println("使用默认的过期时间")
		r.SetDefault(k, v)
		return
	}

	now := time.Now()
	timeDiffer := value.ExpirationTimeStamp - now.Unix()

	// 如果设置了过期时间并且过期时间大于现在将保留原始的过期时间
	if value.IsExpiration && timeDiffer > 0 {
		// fmt.Println("重新计算过期时间")
		// fmt.Println("旧的过期时间", value.ExpirationTimeStamp)
		// fmt.Println("时间限制差", timeDiffer)
		// fmt.Println("新的过期时间", now.Unix()+timeDiffer)
		r.Set(k, v, time.Second*time.Duration(timeDiffer))
	} else {
		// fmt.Println("使用默认的过期时间")
		r.SetDefault(k, v)
	}
}

// 删除 cache
func (r *RedisCacheStruct[T]) Delete(k string) {
	r.Redis.HDel(r.Ctx, r.HashKey, k)
}

// Add() 加入缓存
// func (r *RedisCacheStruct[T]) Add(k string, v T, d time.Duration) {
// 	c.gocahce.Add(k, x, d)
// }

// IncrementInt() 对已存在的key 值自增n
// func (r *RedisCacheStruct[T]) IncrementInt(k string, n int) (num int, err error) {
// 	if err := r.Redis.HIncrBy(r.Ctx, r.HashKey, k, int64(n)).Err(); err != nil {
// 		return num, err
// 	}

// 	if v, ok := r.Get(k); ok {
// 		switch T {
// 		case int:

// 		}
// 		if vint, okint := v.(int); okint {

// 		}
// 	}
// 	return c.gocahce.IncrementInt(k, n)
// }

// ItemCount 获取已存在key的数量
func (r *RedisCacheStruct[T]) ItemCount() (int64, error) {
	if count, err := r.Redis.HLen(r.Ctx, r.HashKey).Result(); err != nil {
		return 0, err
	} else {
		return count, nil
	}
}

// Flush 删除当前已存在的所有key
func (r *RedisCacheStruct[T]) Flush() {
	r.Redis.Del(r.Ctx, r.HashKey)
}

// 定时清理过期验证
func (r *RedisCacheStruct[T]) expirationVerification() {
	ticker := time.NewTicker(r.CleanupInterval)

	for {
		select {
		case <-ticker.C:
			if fields, err := r.Redis.HKeys(r.Ctx, r.HashKey).Result(); err == nil {
				for _, v := range fields {
					// r.Redis.HGet(r.Ctx, r.HashKey, v)
					r.Get(v)
					// fmt.Println("redis定时器", v)
				}
			}
			// case <-j.stop:
			// 	ticker.Stop()
			// 	return
		}
	}
}
