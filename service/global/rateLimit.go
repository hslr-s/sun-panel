package global

import (
	"strconv"
	"sun-panel/lib/cache"
)

type RateLimiter struct {
	Minute cache.Cacher[int]
	Hour   cache.Cacher[int]
}

func (r *RateLimiter) MinuteAddOnce(userId uint) {
	key := "user_" + strconv.Itoa(int(userId))
	times := r.MinuteGet(userId) + 1
	// fmt.Println("fen册数", times)
	r.Minute.SetKeepExpiration(key, times)
	// r.Minute.SetDefault(key, times)
}

func (r *RateLimiter) MinuteGet(userId uint) int {
	if v, ok := r.Minute.Get("user_" + strconv.Itoa(int(userId))); !ok {
		return 0
	} else {
		return v
	}
}

func (r *RateLimiter) HourAddOnce(userId uint) {
	key := "user_" + strconv.Itoa(int(userId))
	times := r.HourGet(userId) + 1
	// fmt.Println("hour册数", times)
	r.Hour.SetKeepExpiration(key, times)
	// r.Hour.SetDefault(key, times)
}

func (r *RateLimiter) HourGet(userId uint) int {
	if v, ok := r.Hour.Get("user_" + strconv.Itoa(int(userId))); !ok {
		return 0
	} else {
		return v
	}
}
