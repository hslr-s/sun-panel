package system

import (
	"fmt"
	"math/rand"
	"sun-panel/api/api_v1/common/apiData/commonApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/lib/cmn/systemSetting"
	"sun-panel/lib/mail"
	"sun-panel/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type registerInfo struct {
	Email        string                               `json:"email"`
	UserName     string                               `json:"userName"`
	Password     string                               `json:"password"`
	Vcode        string                               `json:"vcode"`
	EmailVCode   string                               `json:"emailVCode"`
	VCode        string                               `json:"vCode"`
	Verification commonApiStructs.VerificationRequest `json:"verification"`
	ReferralCode string                               `json:"referralCode"`
}

const EmailCodeCapacity = 1000

type RegisterApi struct{}

// 获取注册验证码
func (l RegisterApi) SendRegisterVcode(c *gin.Context) {
	req := registerInfo{}
	err := c.ShouldBindJSON(&req)
	req.Email = req.UserName
	if err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	errMsg, err := base.ValidateInputStruct(req)
	if err != nil {
		apiReturn.Error(c, errMsg)
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

	// 验证是否开启注册和后缀格式是否正确
	{
		systemSettingInfo := systemSetting.ApplicationSetting{}
		if err := global.SystemSetting.GetValueByInterface("system_application", &systemSettingInfo); err != nil || !systemSettingInfo.Register.OpenRegister {
			apiReturn.Error(c, global.Lang.Get("register.unopened_register"))
			return
		}

		if systemSettingInfo.Register.EmailSuffix != "" && !cmn.VerifyFormat("^.*"+systemSettingInfo.Register.EmailSuffix+"$", req.Email) {
			apiReturn.Error(c, global.Lang.GetWithFields("register.emailSuffix_error", map[string]string{"EmailSuffix": systemSettingInfo.Register.EmailSuffix}))
			return
		}
	}

	// 验证邮箱是否被注册
	{
		userCheck := &models.User{Mail: req.UserName}
		if _, err := userCheck.GetUserInfoByUsername(req.UserName); err == nil && err != gorm.ErrRecordNotFound {
			apiReturn.Error(c, global.Lang.Get("register.mail_exist"))
			return
		}
	}

	emailCode := generateEmailCode()
	count, err := global.VerifyCodeCachePool.ItemCount()
	if err != nil || count >= EmailCodeCapacity {
		global.VerifyCodeCachePool.Flush()
	}
	global.VerifyCodeCachePool.Set(req.Email, emailCode, 0)
	emailInfoConfig := systemSetting.Email{}
	global.SystemSetting.GetValueByInterface("system_email", &emailInfoConfig)
	emailInfo := mail.EmailInfo{
		Username: emailInfoConfig.Mail,
		Password: emailInfoConfig.Password,
		Host:     emailInfoConfig.Host,
		Port:     emailInfoConfig.Port,
	}
	err = mail.SendRegisterEmail(mail.NewEmailer(emailInfo), req.Email, emailCode)
	if err != nil {
		apiReturn.Error(c, global.Lang.Get("mail.send_mail_fail"))
		global.Logger.Errorf("[register] fail to send email to%s", req.UserName)
		return
	}
	apiReturn.Success(c)
}

// 注册提交（开始注册）
func (l *RegisterApi) Commit(c *gin.Context) {
	req := registerInfo{}
	err := c.ShouldBindJSON(&req)
	req.Email = req.UserName
	if err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	errMsg, err := base.ValidateInputStruct(req)
	if err != nil {
		apiReturn.Error(c, errMsg)
		return
	}

	// 验证是否开启注册和后缀格式是否正确
	{
		systemSettingInfo := systemSetting.ApplicationSetting{}
		if err := global.SystemSetting.GetValueByInterface("system_application", &systemSettingInfo); err != nil || !systemSettingInfo.Register.OpenRegister {
			apiReturn.Error(c, global.Lang.Get("register.unopened_register"))
			return
		}

		if systemSettingInfo.Register.EmailSuffix != "" && !cmn.VerifyFormat("^.*"+systemSettingInfo.Register.EmailSuffix+"$", req.Email) {
			apiReturn.Error(c, global.Lang.GetWithFields("register.emailSuffix_error", map[string]string{"EmailSuffix": systemSettingInfo.Register.EmailSuffix}))
			return
		}
	}

	// 验证邮箱是否被注册
	{
		userCheck := &models.User{Mail: req.UserName}
		if _, err := userCheck.GetUserInfoByUsername(req.UserName); err == nil && err != gorm.ErrRecordNotFound {
			apiReturn.Error(c, global.Lang.Get("register.mail_exist"))
			return
		}
	}

	// 验证码验证
	{
		vCode, ok := global.VerifyCodeCachePool.Get(req.Email)
		if !ok {
			apiReturn.Error(c, global.Lang.Get("common.captcha_code_error"))
			//验证码不存在
			return
		}
		if vCode != req.EmailVCode {
			apiReturn.Error(c, global.Lang.Get("common.captcha_code_error"))
			return
			//验证码有误
		}
	}

	// 自动生成用户昵称
	name := "用户" + cmn.BuildRandCode(4, cmn.RAND_CODE_MODE3)

	//验证通过，注册
	user := &models.User{
		Mail:     req.UserName,
		Name:     name,
		Username: req.UserName,
		Password: cmn.PasswordEncryption(req.Password),
		Status:   1,
		Role:     2,
	}
	_, err = user.CreateOne()
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	//删除旧的验证码
	global.VerifyCodeCachePool.Delete(req.Email)

	apiReturn.Success(c)
}

func generateEmailCode() string {
	return fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
}
