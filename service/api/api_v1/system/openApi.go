package system

import (
	"errors"
	"sun-panel/api/api_v1/common/apiData/systemApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/openAPI_v1/common/base"
	"sun-panel/biz"
	"sun-panel/global"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OpenApiApi struct {
}

func (a *OpenApiApi) GetTokenStatusInfo(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	tokenInfo, gormErr := biz.OpenApi.GetUserTokenByUserId(userInfo.ID)
	status := true
	if gormErr != nil {
		if !errors.Is(gormErr, gorm.ErrRecordNotFound) {
			global.Logger.Error(apiReturn.ErrorAndGetLogErrNum(c, gormErr))
			return
		}
		status = false
	}

	apiReturn.SuccessData(c, systemApiStructs.OpenApiGetTokenStatusInfoResp{
		OpenApiTokenInfo: systemApiStructs.OpenApiTokenInfo{
			Token: tokenInfo.Token,
		},
		Status: status,
	})
}

// 刷新token
func (a *OpenApiApi) RefreshToken(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	token, gormErr := biz.OpenApi.RefreshToken(userInfo.ID)

	if gormErr != nil {
		global.Logger.Error(apiReturn.ErrorAndGetLogErrNum(c, gormErr))
		return
	}

	apiReturn.SuccessData(c, systemApiStructs.OpenApiTokenInfo{
		Token: token,
	})
}

func (a *OpenApiApi) DeleteToken(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)

	if err := biz.OpenApi.DeleteToken(userInfo.ID); err != nil {
		global.Logger.Error(apiReturn.ErrorAndGetLogErrNum(c, err))
		return
	}
	apiReturn.Success(c)
}

// 获取路由
// func (a *OpenApiApi) GetRoute(c *gin.Context) {

// 	route := "/openapi/v1/item/create"

// }
