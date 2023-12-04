package base

import (
	"fmt"
	"reflect"
	"strings"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/lib/captcha"
	"sun-panel/lib/cmn"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

type PageLimitVerify struct {
	Page  int64
	Limit int64
}

const (
	VISIT_MODE_LOGIN = iota
	VISIT_MODE_PUBLIC
)

const (
	GIN_GET_VISIT_MODE = "VISIT_MODE"
)

// 验证输入是否有效并返回错误
func validateInputStruct(params interface{}) (errMsg string, err error) {
	var validate = validator.New()
	//通过label标签返回自定义错误内容
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		if label == "" {
			return field.Name
		}
		return label
	})

	// 自定义验证规则，使用 strings.TrimSpace 函数删除前后空格
	validate.RegisterValidation("trimmedRequired", func(fl validator.FieldLevel) bool {
		return strings.TrimSpace(fl.Field().String()) != ""
	})

	if err = validate.Struct(params); err != nil {
		trans := validateTransInit(validate)
		verrs := err.(validator.ValidationErrors)
		// errs := make(map[string]string)
		for _, value := range verrs.Translate(trans) {
			// errs[key[strings.Index(key, ".")+1:]] = value
			errMsg += " " + value
		}
		// fmt.Println(errs)
	}
	return
}

// 验证输入是否有效并返回错误
func ValidateInputStruct(params interface{}) (errMsg string, err error) {
	return validateInputStruct(params)
}

// 数据验证翻译器
func validateTransInit(validate *validator.Validate) ut.Translator {
	// 万能翻译器，保存所有的语言环境和翻译数据
	uni := ut.New(zh.New())
	// 翻译器
	trans, _ := uni.GetTranslator("zh")
	//验证器注册翻译器
	err := zhTranslations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		fmt.Println(err)
	}
	return trans
}

func GetCurrentUserInfo(c *gin.Context) (userInfo models.User, exist bool) {
	if value, exist := c.Get("userInfo"); exist {
		if v, ok := value.(models.User); ok {
			return v, exist
		}
	}
	return
}

// 获取当前访问模式
func GetCurrentVisitMode(c *gin.Context) (visitMode int) {
	if value, exist := c.Get(GIN_GET_VISIT_MODE); exist {
		if v, ok := value.(int); ok {
			return v
		}
	}
	return
}

// 验证器验证
func VerificationCheck(verificationId, vCode string) (errCode int, verificationIdRes string) {

	// 需要进一步验证并返回验证信息
	if verificationId == "" || vCode == "" {
		verificationIdRes = cmn.BuildRandCode(16, cmn.RAND_CODE_MODE1)
		errCode = apiReturn.ERROR_CODE_VERIFICATION_MUST
		return
	}

	// 验证码错误
	if !captcha.CaptchaVerifyHandle(verificationId, vCode) {
		errCode = apiReturn.ERROR_CODE_VERIFICATION_FAIL
		return
	}
	errCode = apiReturn.ERROR_CODE_SUCCESS
	return
}
