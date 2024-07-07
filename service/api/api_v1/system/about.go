package system

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/lib/cmn"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type About struct {
}

type AboutSetAutoCheckNewVersionReq struct {
	biz.VersionAutoSystemConfig
}

func (a *About) Get(c *gin.Context) {
	version := cmn.GetSysVersionInfo()
	apiReturn.SuccessData(c, gin.H{
		"versionName": version.Version,
		"versionCode": version.Version_code,
	})
}

func (a *About) CheckNewVersion(c *gin.Context) {
	type Req struct {
		IsCheckBeta bool
	}

	req := Req{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	var newVersion *biz.LatestVersion

	v, err := biz.ProAuth.CheckNewVersion(biz.VersionTypeRelease)
	if err != nil {
		apiReturn.Error(c, "fail2")
		return
	}

	if v != nil {
		newVersion = v
	}

	if req.IsCheckBeta {
		v, err := biz.ProAuth.CheckNewVersion(biz.VersionTypeBeta)
		if err != nil {
			apiReturn.Error(c, "fail3")
			return
		}

		if v != nil {
			newVersion = v
		}
	}

	apiReturn.SuccessData(c, newVersion)
}

func (a *About) SetAutoCheckNewVersion(c *gin.Context) {
	req := AboutSetAutoCheckNewVersionReq{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if err := biz.Version.SetAutoCheckNewVersionConfig(req.VersionAutoSystemConfig); err != nil {
		global.Logger.Errorln(apiReturn.ErrorAndGetLogErrNum(c, err))
		return
	}

	apiReturn.Success(c)
}

func (a *About) GetAutoCheckNewVersion(c *gin.Context) {
	v, err := biz.Version.GetAutoCheckNewVersionConfig()
	if err != nil {
		global.Logger.Errorln(apiReturn.ErrorAndGetLogErrNum(c, err))
		return
	}
	apiReturn.SuccessData(c, v)
}
