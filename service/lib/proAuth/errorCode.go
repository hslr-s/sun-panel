package proAuth

const (
	ErrNotLogged                   = 1000 // 还未登录
	ErrTokenExpires                = 1001 // token过期
	ErrIncorrectUsernameOrPassword = 1003 // 用户名或密码错误
	ErrAccountDisabled             = 1004 // 账号已停用
	ErrAccountNotExist             = 1006 // 账号不存在
	ErrClientIdNotExist            = 1016 // 客户端不存在
	ErrParameterFormatIncorrect    = 1400 // 参数格式错误
)
