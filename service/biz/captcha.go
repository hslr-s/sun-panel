package biz

import (
	"errors"
	"sun-panel/lib/captcha"

	"github.com/gin-gonic/gin"
)

type CaptchaType struct {
}

func (cpt *CaptchaType) CaptchaVerifyHandle(c *gin.Context, vcode string, clear bool) error {

	var captchaId string
	var err error

	// 获取captchaId
	if captchaId, err = captcha.CaptchaGetIdByCookieHeader(c, "CaptchaId"); err != nil {
		return errors.New("graphical verification code error, please try refreshing the page and try again")
	}

	// 验证码错误
	if !captcha.CaptchaVerifyHandle(captchaId, vcode, clear) {
		captcha.CaptchaVerifyHandle(captchaId, vcode, true) // 此为清除验证码
		return errors.New("graphical verification code error")
	}
	return nil
}
