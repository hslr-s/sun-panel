package apiReturn

import "errors"

const (
	// 特殊
	ErrCodeServerError = 1100

	// 公共
	ErrCodeTokenExpires     = 1000
	ErrCodeParamformatError = 1001
	ErrCodeFailed           = -1

	// 部分接口会出现的
	ErrCodeOnlyNameExisted = 1202
	ErrCodeNoRecord        = 1203
	// ErrCodeSaveIconFileFailed = 1204
)

var ErrorCodeMap = map[int]error{
	ErrCodeServerError: errors.New("server error"),

	ErrCodeTokenExpires:     errors.New("token expires"),
	ErrCodeParamformatError: errors.New("param format error"),

	ErrCodeOnlyNameExisted: errors.New("only name existed"),
	ErrCodeNoRecord:        errors.New("no record"),
	// ErrCodeSaveIconFileFailed: errors.New("failed to save icon file"),
}

func GetError(code int) error {
	if err, ok := ErrorCodeMap[code]; ok {
		return err
	}
	return errors.New("unknown error")
}
