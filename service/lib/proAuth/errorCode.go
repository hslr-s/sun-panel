package proAuth

import "errors"

const (
	ErrNotLogged                   = 1000 // 还未登录
	ErrTokenExpires                = 1001 // token过期
	ErrIncorrectUsernameOrPassword = 1003 // 用户名或密码错误
	ErrAccountDisabled             = 1004 // 账号已停用
	ErrAccountNotExist             = 1006 // 账号不存在
	ErrClientIdNotExist            = 1016 // 客户端不存在
	ErrParameterFormatIncorrect    = 1400 // 参数格式错误
	ErrManualLogin                 = 1500 // 需要手动登录
)

var ErrorCodeMap = map[int]error{
	ErrNotLogged:                   errors.New("not logged in yet"),                 // 还未登录
	ErrTokenExpires:                errors.New("token expires"),                     // 还未登录
	ErrIncorrectUsernameOrPassword: errors.New("incorrect username or password"),    // 用户名或密码错误
	ErrAccountDisabled:             errors.New("account disabled or not activated"), // 账号已停用或未激活
	ErrAccountNotExist:             errors.New("account does not exist"),            // 账号不存在
	ErrClientIdNotExist:            errors.New("clientId does not exist"),           // 客户端不存在
	ErrParameterFormatIncorrect:    errors.New("parameter format error"),            // 参数格式错误
	ErrManualLogin:                 errors.New("need manual login"),                 // 需要手动登录
}

func GetError(code int) error {
	if err, ok := ErrorCodeMap[code]; ok {
		return err
	}
	return errors.New("unknown error")
}

func GetErrorAndDefault(code int, defaultErr error) error {
	if err, ok := ErrorCodeMap[code]; ok {
		return err
	}
	return defaultErr
}
