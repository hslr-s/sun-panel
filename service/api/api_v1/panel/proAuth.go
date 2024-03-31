package panel

import (
	"errors"
	"sun-panel/api/api_v1/common/apiData/panelApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/lib/proAuth"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type ProAuthApi struct {
}

func (p *ProAuthApi) Login(c *gin.Context) {
	// userInfo, _ := base.GetCurrentUserInfo(c)
	req := panelApiStructs.LoginReq{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if err := biz.ProAuth.Login(req.Username, req.Password); err != nil {
		// 账号已停用
		if errors.Is(err, proAuth.GetError(proAuth.ErrIncorrectUsernameOrPassword)) {
			apiReturn.ErrorByCode(c, 1411)
			return
		} else if errors.Is(err, proAuth.GetError(proAuth.ErrAccountDisabled)) {
			apiReturn.ErrorByCode(c, 1412)
			return
		} else if errors.Is(err, proAuth.GetError(proAuth.ErrAccountNotExist)) {
			apiReturn.ErrorByCode(c, 1413)
			return
		} else {
			global.Logger.Error("biz.ProAuth.Login:", err.Error())
			apiReturn.ErrorByCode(c, 1414)
			return
		}
	}

	// 记住账号密码
	biz.ProAuth.SetAccountPassword(req.Username, req.Password)

	apiReturn.Success(c)
}

func (p *ProAuthApi) GetStatus(c *gin.Context) {

	s := biz.ProAuth.GetStatus()

	resp := panelApiStructs.GetStatusResp{
		ProExpiration: s.ProExpiration,
		LoginStatus:   !s.AutoLoginDisable,
	}

	apiReturn.SuccessData(c, resp)
}
