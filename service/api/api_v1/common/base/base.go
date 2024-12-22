package base

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/openAPI_v1/common/base"
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/lib/captcha"
	"sun-panel/lib/cmn"
	"sun-panel/models"
	"time"

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

var (
	ErrUploadExtensionNameNotAllowed = errors.New("extension not allowed") // "上传文件类型不允许"
	ErrUploadExceedMaxSize           = errors.New("exceeds maximum size")  // "上传文件超出最大尺寸"
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

// // 保存图片并将记录保存到数据库中
// func UploadImageAndSaveImage(c *gin.Context, f *multipart.FileHeader, source_path string, fileExt string, userID uint, fileType int) (url string) {

// 	fileName := cmn.Md5(fmt.Sprintf("%s%s", f.Filename, time.Now().String()))
// 	fildDir := fmt.Sprintf("%s/%d/%d/%d/", source_path, time.Now().Year(), time.Now().Month(), time.Now().Day())
// 	isExist, _ := cmn.PathExists(fildDir)
// 	if !isExist {
// 		os.MkdirAll(fildDir, os.ModePerm)
// 	}
// 	filepath := fmt.Sprintf("%s%s%s", fildDir, fileName, fileExt)
// 	c.SaveUploadedFile(f, filepath)

// 	pureFilePath := filepath[len(source_path):]
// 	// 向数据库添加记录
// 	mFile := models.File{}
// 	mFile.AddFile(userID, f.Filename, fileExt, pureFilePath, fileType)
// 	return global.UPLOAD_ROUTE + pureFilePath
// 	// apiReturn.SuccessData(c, gin.H{
// 	// 	"imageUrl": global.UPLOAD_ROUTE + pureFilePath,
// 	// })
// }

// 保存图片并将记录保存到数据库中
func UploadFileAndSave(c *gin.Context, agreeExtNames []string, maxSize int64) (fileInfo models.File, err error) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	configUpload := global.Config.GetValueString("base", "source_path")

	// 文件类型
	fileTypeStr := c.PostForm("type")
	fileType, _ := strconv.Atoi(fileTypeStr)

	f, err := c.FormFile("imgfile")
	if err != nil {
		return
	} else {
		fileExt := strings.ToLower(path.Ext(f.Filename))

		if !cmn.InArray(agreeExtNames, fileExt) {
			err = ErrUploadExtensionNameNotAllowed
			return
		}

		if maxSize != 0 && f.Size > maxSize {
			err = ErrUploadExceedMaxSize
			return
		}

		fileName := cmn.Md5(fmt.Sprintf("%s%s", f.Filename, time.Now().String()))
		fildDir := fmt.Sprintf("%s/%d/%d/%d/", configUpload, time.Now().Year(), time.Now().Month(), time.Now().Day())
		isExist, _ := cmn.PathExists(fildDir)
		if !isExist {
			os.MkdirAll(fildDir, os.ModePerm)
		}
		filepath := fmt.Sprintf("%s%s%s", fildDir, fileName, fileExt)
		c.SaveUploadedFile(f, filepath)

		pureFilePath := filepath[len(configUpload):]
		// 向数据库添加记录
		fileInfo = models.File{
			UserId:   userInfo.ID,
			FileName: f.Filename,
			Src:      pureFilePath,
			Type:     fileType,
			Ext:      fileExt,
		}

		err = global.Db.Create(&fileInfo).Error
	}

	return
}

// 下载图片并将记录保存到数据库中
func DownloadFileAndSave(fileUrl string, userId uint) (fileInfo models.File, err error) {
	return biz.File.DownloadFileAndSaveDb(global.Db, fileUrl, userId)
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
	if !captcha.CaptchaVerifyHandle(verificationId, vCode, true) {
		errCode = apiReturn.ERROR_CODE_VERIFICATION_FAIL
		return
	}
	errCode = apiReturn.ERROR_CODE_SUCCESS
	return
}

func AnyToJsonStr(d interface{}) string {
	b, err := json.Marshal(d)
	bs := string(b)
	if err != nil {
		bs = "{}"
	}
	return bs
}

func JsonStrToAny(s string, value interface{}) error {
	return json.Unmarshal([]byte(s), value)
}

func JsonStrToMap(s string) map[string]interface{} {
	data := map[string]any{}
	JsonStrToAny(s, &data)
	return data
}

func IntToBool(v int) bool {
	return v == 1
}

func BoolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
