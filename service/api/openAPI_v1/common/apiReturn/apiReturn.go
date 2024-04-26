package apiReturn

import (
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"time"

	"github.com/gin-gonic/gin"
)

const ERROR_CODE_SUCCESS = 0 // 错误码 无任何错误

func ApiReturn(ctx *gin.Context, code int, msg string, data interface{}) {
	returnData := map[string]interface{}{
		"code": code,
		"msg":  msg,
	}
	if data != nil {
		returnData["data"] = data
	}
	ctx.JSON(200, returnData)
}

// 返回成功
func SuccessData(ctx *gin.Context, data interface{}) {
	ApiReturn(ctx, 0, "OK", data)
}

// 返回列表
func SuccessListData(ctx *gin.Context, list interface{}, count int64) {
	ApiReturn(ctx, 0, "OK", gin.H{
		"list":  list,
		"count": count,
	})
}

// 返回成功，没有data数据
func Success(ctx *gin.Context) {
	ApiReturn(ctx, 0, "OK", nil)
}

// 返回错误 需要个性化定义的错误|带返回数据的错误
func ErrorCode(ctx *gin.Context, code int, errMsg string, data interface{}) {
	ApiReturn(ctx, code, errMsg, data)
}

// 返回错误 普通提示错误
func Error(ctx *gin.Context, errMsg string) {
	ErrorCode(ctx, -1, errMsg, nil)
}

// 返回错误 参数错误
func ErrorParamFomat(ctx *gin.Context, errMsg string) {
	ErrorByCodeAndMsg(ctx, ErrCodeParamformatError, GetError(ErrCodeParamformatError).Error()+":"+errMsg)
}

// // 返回错误 数据库
// func ErrorDatabase(ctx *gin.Context, errMsg string) {
// 	// Error(ctx, global.Lang.GetAndInsert("common.db_error", "[", errMsg, "]"))
// 	ErrorByCodeAndMsg(ctx, 1200, errMsg)

// }

// 返回错误 数据记录未找到
func ErrorDataNotFound(ctx *gin.Context) {
	// ErrorCode(ctx, ERROR_CODE_DATA_RECORD_NOT_FOUND, "未找到数据记录", nil)
	ErrorByCode(ctx, -1)
}

func ErrorByCode(ctx *gin.Context, code int) {
	msg := "Server error"
	if v, ok := GetErrorMsgByCode(code); ok {
		msg = v
	}
	ErrorCode(ctx, code, msg, nil)
}

// 返回错误编号，并记录到日志中，对请求方不可见具体错误
func ErrorWithNum(ctx *gin.Context, err error) {
	errNum := time.Now().Format("20060102150405") + cmn.BuildRandCode(5, "0123456789")

	returnData := map[string]interface{}{
		"code": ErrCodeServerError,
		"msg":  GetError(ErrCodeServerError).Error() + ": ErrNum:" + errNum,
	}
	global.Logger.Error("ErrNum:"+errNum, "-", err)
	ctx.JSON(200, returnData)
}

// 使用错误码的错误并附加错误信息
func ErrorByCodeAndMsg(ctx *gin.Context, code int, msg string) {
	defalurMsg := "Server error"
	if v, ok := GetErrorMsgByCode(code); ok {
		msg = v
	}
	ErrorCode(ctx, code, defalurMsg+"["+msg+"]", nil)
}

func GetErrorMsgByCode(code int) (string, bool) {
	if v, ok := ErrorCodeMap[code]; ok {
		return v.Error(), true
	} else {
		return "", false
	}
}
