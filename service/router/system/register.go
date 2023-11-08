package system

import (
	"sun-panel/api/api_v1"

	"github.com/gin-gonic/gin"
)

func InitRegister(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.RegisterApi

	router.POST("/register/sendRegisterVcode", api.SendRegisterVcode)
	router.POST("/register/commit", api.Commit)
}
