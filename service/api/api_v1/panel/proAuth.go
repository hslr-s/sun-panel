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
		// 账号已停用，避免被破解统一密码或者账号错误
		if errors.Is(err, proAuth.GetError(proAuth.ErrIncorrectUsernameOrPassword)) {
			apiReturn.ErrorByCode(c, apiReturn.ErrCodeIncorrectUsernameOrPassword)
			return
		} else if errors.Is(err, proAuth.GetError(proAuth.ErrAccountDisabled)) {
			apiReturn.ErrorByCode(c, apiReturn.ErrCodeIncorrectUsernameOrPassword)
			return
		} else if errors.Is(err, proAuth.GetError(proAuth.ErrAccountNotExist)) {
			apiReturn.ErrorByCode(c, apiReturn.ErrCodeIncorrectUsernameOrPassword)
			return
		} else {
			global.Logger.Error("biz.ProAuth.Login:", err.Error())
			apiReturn.Error(c, err.Error())
			return
		}
	}

	// 记住账号密码
	biz.ProAuth.SetAccountPassword(req.Username, req.Password)

	apiReturn.Success(c)
}

func (p *ProAuthApi) GetStatus(c *gin.Context) {

	s := biz.ProAuth.GetStatus()

	shortAuthData := &biz.ProAuthDataShortTerm{}
	if s.ProAuthDataShortTerm != nil {
		shortAuthData = s.ProAuthDataShortTerm
	}

	resp := panelApiStructs.GetStatusResp{
		ProExpiration: shortAuthData.ProExpiration,
		LoginStatus:   s.LoginStatus,
		Username:      shortAuthData.UserInfo.Username,
		NickName:      shortAuthData.UserInfo.Name,
	}

	apiReturn.SuccessData(c, resp)
}

// 刷新状态
func (p *ProAuthApi) RefreshStatus(c *gin.Context) {
	err := biz.ProAuth.RefreshInfo()
	if err != nil {
		if errors.Is(err, proAuth.GetError(proAuth.ErrTokenExpires)) {
			apiReturn.ErrorByCode(c, apiReturn.ErrCodeTokenExpired)
			return
		} else {
			global.Logger.Error("biz.ProAuth.Login:", err.Error())
			apiReturn.Error(c, err.Error())
			return
		}
	}
	apiReturn.Success(c)
}

// 退出账号
func (p *ProAuthApi) Logout(c *gin.Context) {
	biz.ProAuth.Reset()
	apiReturn.Success(c)
}

// 设置隐藏的PRO徽章的状态
func (p *ProAuthApi) SetHideProBadgeStatus(c *gin.Context) {
	type Req struct {
		Status bool
	}

	req := Req{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	biz.ProAuth.SetHideProBadgeStatus(req.Status)
	apiReturn.Success(c)
}
