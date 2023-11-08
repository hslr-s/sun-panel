package rateLimit

import (
	"errors"
	"sun-panel/global"
)

const (
	ERROR_RATE_EXCEED_MINUTE = "minute exceed" // 分钟速率超出限制
	ERROR_RATE_EXCEED_HOUR   = "minute hour"   // 小时速率超出限制
)

// 获取用户套餐的速率 此处正常根据用户套餐设定获取-暂时写死
func GetUserPackageRate(userId uint) (minuteRate, hourRate int) {
	return 10, 200
}

func CheckRateLimit(userId uint) error {
	minuteRate, hourRate := GetUserPackageRate(userId)
	if minuteRate != 0 && minuteRate <= global.RateLimit.MinuteGet(userId) {
		return errors.New(ERROR_RATE_EXCEED_MINUTE)
	}

	if hourRate != 0 && hourRate <= global.RateLimit.HourGet(userId) {
		return errors.New(ERROR_RATE_EXCEED_HOUR)
	}

	return nil
}

// 速率+1次 同时增加小时和分钟的次数
func AddOnceRate(userId uint) error {
	global.RateLimit.MinuteAddOnce(userId)
	global.RateLimit.HourAddOnce(userId)
	return nil
}
