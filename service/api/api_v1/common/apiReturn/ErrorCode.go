package apiReturn

import "errors"

const (
	ErrCodeNotLoggedIn                   = 1000
	ErrCodeTokenExpired                  = 1001
	ErrCodeCaptchaError                  = 1002
	ErrCodeIncorrectUsernameOrPassword   = 1003
	ErrCodeAccountDisabledOrNotActivated = 1004
	ErrCodeNoCurrentPermission           = 1005
	ErrCodeAccountDoesNotExist           = 1006
	ErrCodeOldPasswordError              = 1007
	ErrCodeNoPROAuthorization            = 1008

	ErrCodeDatabaseError        = 1200
	ErrCodePleaseKeepAtLeastOne = 1201
	ErrCodeNoDataRecordFound    = 1202

	ErrCodeUploadFailed          = 1300
	ErrCodeUnsupportedFileFormat = 1301

	ErrCodeParameterFormatError                  = 1400
	ErrCodeProAuthIncorrectUsernameOrPassword    = 1411
	ErrCodeProAuthAccountDisabledOrNotActivated  = 1412
	ErrCodeProAuthAccountDoesNotExist            = 1413
	ErrCodeProAuthServerError                    = 1414
	ErrCodeCustomJSAndCSSFilePreservationFailure = 1430
)

var ErrorCodeMap = map[int]error{
	ErrCodeNotLoggedIn:                   errors.New("not logged in yet"),                 // 还未登录
	ErrCodeTokenExpired:                  errors.New("token expired"),                     // token过期
	ErrCodeCaptchaError:                  errors.New("captcha Error"),                     // 验证码错误
	ErrCodeIncorrectUsernameOrPassword:   errors.New("incorrect username or password"),    // 用户名或密码错误
	ErrCodeAccountDisabledOrNotActivated: errors.New("account disabled or not activated"), // 账号已停用或未激活
	ErrCodeNoCurrentPermission:           errors.New("no permission"),                     // 无权限
	ErrCodeAccountDoesNotExist:           errors.New("account does not exist"),            // 账号不存在
	ErrCodeOldPasswordError:              errors.New("old password error"),                // 旧密码不正确
	ErrCodeNoPROAuthorization:            errors.New("no PRO authorization"),              // 没有PRO授权

	// 数据类
	ErrCodeDatabaseError:        errors.New("database error"),           // 数据库错误
	ErrCodePleaseKeepAtLeastOne: errors.New("please keep at least one"), // 请至少保留一个
	ErrCodeNoDataRecordFound:    errors.New("no data record found"),     // 未找到数据记录

	ErrCodeUploadFailed:          errors.New("upload failed"),           // 上传失败
	ErrCodeUnsupportedFileFormat: errors.New("unsupported file format"), // 不被支持的格式文件

	ErrCodeParameterFormatError: errors.New("parameter format error"), // 参数格式错误

	ErrCodeProAuthIncorrectUsernameOrPassword:   errors.New("incorrect username or password"),    // proAuth:账号密码错误
	ErrCodeProAuthAccountDisabledOrNotActivated: errors.New("account disabled or not activated"), // proAuth:账号已停用
	ErrCodeProAuthAccountDoesNotExist:           errors.New("account does not exist"),            // proAuth:账号不存在
	ErrCodeProAuthServerError:                   errors.New("server error"),                      // proAuth:服务器错误

	ErrCodeCustomJSAndCSSFilePreservationFailure: errors.New("custom JS and CSS file preservation failure"),
}

func GetError(code int) error {
	if err, ok := ErrorCodeMap[code]; ok {
		return err
	}
	return errors.New("unknown error")
}
