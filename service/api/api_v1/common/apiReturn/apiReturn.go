package apiReturn

import (
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"time"

	"github.com/gin-gonic/gin"
)

const ERROR_CODE_SUCCESS = 0 // 错误码 无任何错误

const (
	// 验证器类

	ERROR_CODE_VERIFICATION_MUST = 1101 // 错误码 验证器类：必须需要验证或者验证数据为空
	ERROR_CODE_VERIFICATION_FAIL = 1102 // 错误码 验证器类：验证失败，验证错误

	// 数据类

	ERROR_CODE_DATA_DATABASE         = 1200 // 错误码 数据类：数据库报错
	ERROR_CODE_DATA_RECORD_NOT_FOUND = 1202 // 错误码 数据类：数据记录未找到
)

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

// 返回列表数据
func ListData(ctx *gin.Context, list interface{}, count int64) {
	data := map[string]interface{}{
		"list":  list,
		"count": count,
	}
	ApiReturn(ctx, 0, "OK", data)
}

func ErrorWithLogErrNum(ctx *gin.Context, err error) {
	errNum := time.Now().Format("20060102150405") + cmn.BuildRandCode(5, "0123456789")

	returnData := map[string]interface{}{
		"code": -1,
		"msg":  "ErrNum:" + errNum,
	}
	// global.Logger.Errorln("ErrNum:"+errNum, "-", err)
	global.Logger.Errorw("ErrNum:"+errNum, "skip", 1)
	ctx.JSON(200, returnData)
}

// 错误 并获取错误码详情
// 搭配
// global.Logger.Error(apiReturn.ErrorAndGetLogErrNum(c, err))
func ErrorAndGetLogErrNum(ctx *gin.Context, err error) string {
	errNum := time.Now().Format("20060102150405") + cmn.BuildRandCode(5, "0123456789")

	returnData := map[string]interface{}{
		"code": -1,
		"msg":  "ErrNum:" + errNum,
	}

	ctx.JSON(200, returnData)
	return "ErrNum:" + errNum + " | " + err.Error()
}

// 返回错误 需要个性化定义的错误|带返回数据的错误
func ErrorCode(ctx *gin.Context, code int, errMsg string, data interface{}) {
	ApiReturn(ctx, code, errMsg, data)
}

// 返回错误 普通提示错误
func Error(ctx *gin.Context, errMsg string) {
	ErrorCode(ctx, -1, errMsg, nil)
}

// 返回错误 需要个性化定义的错误|带返回数据的错误
func ErrorNoAccess(ctx *gin.Context) {
	ErrorByCode(ctx, ErrCodeNoCurrentPermission)
}

// 返回错误 参数错误
func ErrorParamFomat(ctx *gin.Context, errMsg string) {
	Error(ctx, global.Lang.GetAndInsert("common.api_error_param_format", "[", errMsg, "]"))
	// Error(ctx, "参数错误")
}

// // 返回错误 数据库
func ErrorDatabase(ctx *gin.Context, errMsg string) {
	// Error(ctx, global.Lang.GetAndInsert("common.db_error", "[", errMsg, "]"))
	ErrorByCodeAndMsg(ctx, 1200, errMsg)

}

// 返回错误 数据记录未找到
func ErrorDataNotFound(ctx *gin.Context) {
	// ErrorCode(ctx, ERROR_CODE_DATA_RECORD_NOT_FOUND, "未找到数据记录", nil)
	ErrorByCode(ctx, ErrCodeNoDataRecordFound)
}

func ErrorByCode(ctx *gin.Context, code int) {
	msg := "Server error"
	if v, ok := GetErrorMsgByCode(code); ok {
		msg = v
	}
	ErrorCode(ctx, code, msg, nil)
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

// 返回错误 需要个性化定义的错误|带返回数据的错误
// func ErrorNoAccess(ctx *gin.Context) {
// 	ErrorCode(ctx, 1005, global.Lang.Get("common.no_access"), nil)
// }

// // 返回错误 参数错误
// func ErrorParamFomat(ctx *gin.Context, errMsg string) {
// 	Error(ctx, global.Lang.GetAndInsert("common.api_error_param_format", "[", errMsg, "]"))
// }

// // 返回错误 数据库
// func ErrorDatabase(ctx *gin.Context, errMsg string) {
// 	Error(ctx, global.Lang.GetAndInsert("common.db_error", "[", errMsg, "]"))
// }
