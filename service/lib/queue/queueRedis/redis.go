package queueRedis

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type Pool struct {
	Ctx   context.Context
	Redis *redis.Client
	Name  string
}

func New(redisDb *redis.Client, name string) *Pool {
	ctx := context.Background()
	return &Pool{
		Ctx:   ctx,
		Redis: redisDb,
		Name:  name,
	}
}

func (k *Pool) LPush(value ...interface{}) error {
	var values []interface{}
	for _, v := range value {
		values = append(values, k.encode(v))
	}

	return k.Redis.LPush(k.Ctx, k.Name, values...).Err()
}

func (k *Pool) RPush(value ...interface{}) error {
	var values []interface{}
	for _, v := range value {
		values = append(values, k.encode(v))
	}
	return k.Redis.RPush(k.Ctx, k.Name, values...).Err()
}

func (k *Pool) Delete(value interface{}) error {
	return k.Redis.LRem(k.Ctx, k.Name, 0, k.encode(value)).Err()
}

// 取出值
func (k *Pool) GetByIndex(index int64, v interface{}) error {
	if d, err := k.Redis.LIndex(k.Ctx, k.Name, index).Result(); err != nil {
		return err
	} else {
		k.decode(d, v)
		return nil
	}
}

// 左-取出并删除
func (k *Pool) LPop(v interface{}) error {
	if d, err := k.Redis.LPop(k.Ctx, k.Name).Result(); err != nil {
		return err
	} else {
		k.decode(d, v)
		return nil
	}
}

// 右-取出并删除
func (k *Pool) RPop(v interface{}) error {
	if d, err := k.Redis.RPop(k.Ctx, k.Name).Result(); err != nil {
		return err
	} else {
		k.decode(d, v)
		return nil
	}
}

func (k *Pool) encode(value any) string {
	data, err := json.Marshal(value)
	if err != nil {
		return "{}"
	}
	return string(data)
}

func (k *Pool) decode(v string, value interface{}) {
	err := json.Unmarshal([]byte(v), value)
	_ = err
}

func (k *Pool) Length() (int64, error) {
	r := k.Redis.LLen(k.Ctx, k.Name)
	return r.Result()
}

func (k *Pool) Flush() error {
	return k.Redis.Del(k.Ctx, k.Name).Err()
}
