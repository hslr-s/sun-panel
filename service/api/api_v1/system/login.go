package system

import (
	"strconv"
	"sun-panel/api/api_v1/common/apiData/commonApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/lib/captcha"
	"sun-panel/lib/cmn"
	"sun-panel/lib/cmn/systemSetting"
	"sun-panel/lib/mail"
	"sun-panel/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LoginApi struct {
}

// 登录输入验证
type LoginLoginVerify struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,max=50"`
	VCode    string `json:"vcode" validate:"max=6"`
	Email    string `json:"email"`
}

// @Summary 登录账号
// @Accept application/json
// @Produce application/json
// @Param LoginLoginVerify body LoginLoginVerify true "登陆验证信息"
// @Tags user
// @Router /login [post]
func (l LoginApi) Login(c *gin.Context) {
	param := LoginLoginVerify{}
	if err := c.ShouldBindJSON(&param); err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format"))
		return
	}

	if errMsg, err := base.ValidateInputStruct(param); err != nil {
		apiReturn.Error(c, errMsg)
		return
	}

	settings := systemSetting.ApplicationSetting{}
	global.SystemSetting.GetValueByInterface("system_application", &settings)

	// 验证验证码
	if settings.Login.LoginCaptcha {
		var captchaId string
		var err error

		// 获取captchaId
		if captchaId, err = captcha.CaptchaGetIdByCookieHeader(c, "CaptchaId"); err != nil {
			apiReturn.Error(c, global.Lang.Get("login.err_captcha_check_fail"))
			return
		}

		// 验证码错误
		if !captcha.CaptchaVerifyHandle(captchaId, param.VCode) {
			apiReturn.Error(c, global.Lang.Get("captcha.api_captcha_fail"))
			return
		}
	}

	mUser := models.User{}
	var (
		err  error
		info models.User
	)
	bToken := ""
	if info, err = mUser.GetUserInfoByUsernameAndPassword(param.Username, cmn.PasswordEncryption(param.Password)); err != nil {
		// 未找到记录 账号或密码错误
		if err == gorm.ErrRecordNotFound {
			apiReturn.Error(c, global.Lang.Get("login.err_username_password"))
			return
		} else {
			// 未知错误
			apiReturn.Error(c, err.Error())
			return
		}

	}

	// 停用或未激活
	if info.Status != 1 {
		apiReturn.Error(c, global.Lang.Get("login.err_username_deactivation"))
		return
	}

	bToken = info.Token
	if info.Token == "" {
		// 生成token
		buildTokenOver := false
		for !buildTokenOver {
			bToken = cmn.BuildRandCode(32, cmn.RAND_CODE_MODE2)
			if _, err := mUser.GetUserInfoByToken(bToken); err != nil {
				// 保存token
				mUser.UpdateUserInfoByUserId(info.ID, map[string]interface{}{
					"token": bToken,
				})
				buildTokenOver = true
			}
		}
		info.Token = bToken
	}
	info.Password = ""
	info.ReferralCode = ""

	// global.UserToken.SetDefault(bToken, info)
	cToken := uuid.NewString() + "-" + cmn.Md5(cmn.Md5("userId"+strconv.Itoa(int(info.ID))))
	global.CUserToken.SetDefault(cToken, bToken)
	global.Logger.Debug("token:", cToken, "|", bToken)
	global.Logger.Debug(global.CUserToken.Get(cToken))

	// 设置当前用户信息
	c.Set("userInfo", info)
	info.Token = cToken // 重要 采用cToken,隐藏真实token
	apiReturn.SuccessData(c, info)
}

// 安全退出
func (l *LoginApi) Logout(c *gin.Context) {
	// userInfo, _ := base.GetCurrentUserInfo(c)
	cToken := c.GetHeader("token")
	global.CUserToken.Delete(cToken)
	apiReturn.Success(c)
}

// 获取重置密码的验证码
func (l *LoginApi) SendResetPasswordVCode(c *gin.Context) {
	type ResstRequest struct {
		LoginLoginVerify
		Verification commonApiStructs.VerificationRequest `json:"verification"`
	}
	req := ResstRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format"))
		return
	}

	// 验证码验证
	{
		errCode, verifcationId := base.VerificationCheck(req.Verification.CodeID, req.Verification.VCode)
		if errCode != apiReturn.ERROR_CODE_SUCCESS {
			apiReturn.ErrorVerification(c, errCode, verifcationId)
			return
		}
	}

	emailVCode := cmn.BuildRandCode(6, cmn.RAND_CODE_MODE2)
	global.VerifyCodeCachePool.Set(req.Email, emailVCode, 10*time.Minute)

	userCheck := &models.User{Mail: req.Email}
	userInfo := userCheck.GetUserInfoByMail()
	if userInfo == nil {
		apiReturn.Error(c, "账号不存在")
		return
	}
	emailInfoConfig := systemSetting.Email{}
	global.SystemSetting.GetValueByInterface("system_email", &emailInfoConfig)
	emailInfo := mail.EmailInfo{
		Username: emailInfoConfig.Mail,
		Password: emailInfoConfig.Password,
		Host:     emailInfoConfig.Host,
		Port:     emailInfoConfig.Port,
	}
	if err := mail.SendResetPasswordVCode(mail.NewEmailer(emailInfo), req.Email, emailVCode); err != nil {
		apiReturn.Error(c, err.Error())
		return
	}

	apiReturn.Success(c)

}

// 使用邮箱验证码重置密码
func (l *LoginApi) ResetPasswordByVCode(c *gin.Context) {
	req := registerInfo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format"))
		return
	}

	userCheck := &models.User{Mail: req.Email}
	userInfo := userCheck.GetUserInfoByMail()
	if userInfo == nil {
		apiReturn.Error(c, "账号不存在")
		return
	}

	// 校验验证码
	{
		if emailVCode, ok := global.VerifyCodeCachePool.Get(req.Email); !ok || req.EmailVCode != emailVCode {
			apiReturn.Error(c, global.Lang.Get("common.captcha_code_error"))
			return
		}
		global.VerifyCodeCachePool.Delete(req.Email)
	}

	updateData := map[string]interface{}{
		"password": cmn.PasswordEncryption(req.Password),
		"token":    "",
	}
	global.UserToken.Delete(userInfo.Token) // 更新用户信息
	if err := userInfo.UpdateUserInfoByUserId(userInfo.ID, updateData); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.Success(c)

}
