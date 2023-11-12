package apiReturn

import "github.com/gin-gonic/gin"

// 返回错误docker 连接失败错误
func ErrorDockerClientFail(ctx *gin.Context, content string) {
	ErrorCode(ctx, -1, "Docker ERROR- connection failed : ["+content+"]", nil)
}

func ErrorDocker(ctx *gin.Context, content string) {
	ErrorCode(ctx, -1, "Docker ERROR- ["+content+"]", nil)
}
