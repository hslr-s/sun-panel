package system

import (
	"encoding/base64"
	"strconv"
	"sun-panel/lib/captcha"
	"sun-panel/lib/cmn"

	"github.com/gin-gonic/gin"
)

type CaptchaApi struct {
	ErrMsg string // 错误信息
}

// 获取图像
func (c *CaptchaApi) GetImage(ctx *gin.Context) {
	key := cmn.BuildRandCode(16, cmn.RAND_CODE_MODE2)
	width, _ := strconv.Atoi(ctx.Param("width"))
	height, _ := strconv.Atoi(ctx.Param("height"))
	if width == 0 || width > 500 {
		width = 120
	}
	if height == 0 || height > 500 {
		height = 44
	}
	// 设置网页验证码的cookie
	ctx.SetCookie("CaptchaId", key, 3600, "/", "", false, false)
	base64Str := captcha.GenerateCaptchaHandler(key, width, height)
	_ = base64Str
	// base64 字符串一般会包含头部 data:image/xxx;base64, 需要去除
	baseImg, _ := base64.StdEncoding.DecodeString(base64Str[22:])
	_, _ = ctx.Writer.WriteString(string(baseImg))
}

// 获取图像根据验证器id，id从地址栏获取
func (c *CaptchaApi) GetImageByCaptchaId(ctx *gin.Context) {
	// key := cmn.BuildRandCode(16, cmn.RAND_CODE_MODE2)
	width, _ := strconv.Atoi(ctx.Param("width"))
	height, _ := strconv.Atoi(ctx.Param("height"))
	captchaId := ctx.Param("captchaId")
	if width == 0 || width > 500 {
		width = 120
	}
	if height == 0 || height > 500 {
		height = 44
	}
	// 设置网页验证码的cookie
	base64Str := captcha.GenerateCaptchaHandler(captchaId, width, height)
	_ = base64Str
	// base64 字符串一般会包含头部 data:image/xxx;base64, 需要去除
	baseImg, _ := base64.StdEncoding.DecodeString(base64Str[22:])
	_, _ = ctx.Writer.WriteString(string(baseImg))
}

func (c *CaptchaApi) CheckVCode(id, vcode string) {
	// Captcha.Store = base64Captcha.DefaultMemStore
	// if store.Verify(id, vcode, true) {
	// 	body = map[string]interface{}{"code": 1001, "msg": "ok"}
	// }
	// w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// json.NewEncoder(w).Encode(body)
}
