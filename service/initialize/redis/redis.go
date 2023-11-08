package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Options struct {
	Addr     string // localhost:6379
	Password string // 没有密码，默认值
	DB       int    // 默认DB 0
}

func InitRedis(options Options) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     options.Addr,
		Password: options.Password,
		DB:       options.DB,
	})

	// 验证连接是否成功
	ctx := context.Background()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, err
	}
	return rdb, nil
}
