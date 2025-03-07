package biz

import (
	"encoding/json"
	"errors"
	"net"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/lib/proAuth"
	"sun-panel/lib/proAuth/AES"
	"time"

	"github.com/google/uuid"
)

type ProAuthType struct {
	ProAuthDataShortTerm *ProAuthDataShortTerm // 短期授权数据，为nil，视为未登录
	AutoLoginDisable     bool                  // 为true时需要手动登录
	LoginStatus          bool                  // 登录状态
	PingInterval         time.Duration

	// ProExpiration        *time.Time
	// UserInfo             proAuth.UserInfo
	// Token  string
	// Status bool // 登录状态
}

type ProAuthAccountPassword struct {
	Username string
	Password string
}

type ProAuthDataShortTerm struct {
	// ProAuthAccountPassword
	ShortTermExpiration *time.Time // 短期授权
	ProExpiration       *time.Time // 最终PRO授权
	Token               string
	UserInfo            proAuth.UserInfo // 昵称
}

type LatestVersion struct {
	ReleaseTime time.Time           `json:"releaseTime"`
	Type        proAuth.VersionType `json:"type"`
	Version     string              `json:"version"`
	DownloadURL string              `json:"downloadURL"`
	PageUrl     string              `json:"pageUrl"`
}

type VersionType string

const (
	VersionTypeRelease VersionType = "release"
	VersionTypeBeta    VersionType = "beta"
	// VersionTypeAlpha   VersionType = "alpha"
	// VersionTypeRC      VersionType = "rc"
)

var ErrAccountPasswordIncomplete = errors.New("ErrAccountPasswordIncomplete")

const ProAuthTypeSystemSettingClientIdKey = "clientId"

// 短期授权的key
const ProAuthTypeSystemSettingProAuthDataShortTermKey = "pro_auth_data_short_term"

// 将废弃 (1.4.0-beta24-04-12添加的此项)
const ProAuthTypeSystemSettingAccountPasswordKey = "ProAuthAccountPassword"

var AuthUrl = map[string]string{
	"default": "https://auth.sun-panel.top",
	"spare":   "https://auth.sp.enianteam.com", // 备用验证地址
}

func NewProAuth() *ProAuthType {
	// 此配置项为隐藏的默认为正式服务器，该字段目前仅支持配置 default 和 spare 和 加密的url地址，不再支持自定义明文地址 更新日期：2025-03-07
	cfgAuthHost := global.Config.GetValueStringOrDefault("pro_auth", "host")

	proAuth.Config.AESKey = "wAOhhAqbBxgjNOlcEAGvwawvUZmWDLkN" // 绝密不可泄露
	hostAesKey := "60l1o651mL3WVcbWW0NJLTLav0kpu1vT"           // 绝密不可泄露，用于解密host地址
	proAuth.Config.TimeDifference = 604800                     // 允许的时间差为7天

	useHostFlag := cfgAuthHost
	if cfgAuthHost == "" || cfgAuthHost == "default" {
		proAuth.Config.Url = AuthUrl["default"]
		useHostFlag = "default"
	} else if cfgAuthHost == "spare" {
		proAuth.Config.Url = AuthUrl["spare"]
	} else {
		// 默认认为是一个加密的网址
		url, err := AES.Decrypt(hostAesKey, cfgAuthHost)
		if err == nil {
			proAuth.Config.Url = url
			// global.Logger.Debugln("auth_host host_address:", url)
		} else {
			global.Logger.Errorln("host_address: Unable to parse :", err)
			proAuth.Config.Url = AuthUrl["default"]
			useHostFlag = "default"
		}
	}

	global.Logger.Debugln("auth_host:", proAuth.Config.Url)
	global.Logger.Infoln("Use PRO Auth Host:", useHostFlag)

	p := ProAuthType{
		PingInterval:     6 * time.Hour, // ping执行的时间间隔
		AutoLoginDisable: false,         // 刚启动允许自动登录
		LoginStatus:      false,
	}
	// debug模式连接本地授权平台
	if global.RUNCODE == "debug" {
		p.PingInterval = 30 * time.Second
		// proAuth.Config.IsDebug = true // 明文传输
	}

	return &p
}

func (p *ProAuthType) Start() {

	// 读取数据库/缓存中的token等短期授权的数据
	if v := p.GetAuthDataShortTerm(); v != nil {
		p.LoginStatus = true
		p.ProAuthDataShortTerm = v
	}

	// p.Login("95302870@qq.com", "12345678")
	go p.Ping()
	// p.Register()
	// p.SetAccountPassword("95302870@qq.com", "12345678")
	// p.AutoLogin()
	go p.timer(p.PingInterval, p.timerFunc)
}

func (p *ProAuthType) GetLocalLanIpMAC() (ip, macAddr string, resErr error) {
	errString := ""
	// 获取本地非环回地址
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		resErr = err
		return
	}
	// 遍历所有地址
	for _, address := range addrs {
		// 检查是否是IP网络地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				// 获取与IP地址关联的接口
				interfaces, err := net.Interfaces()
				if err != nil {
					errString = err.Error()
					break
				}
				// 遍历所有接口
				for _, inter := range interfaces {
					// 找到对应的接口
					if (inter.Flags & net.FlagUp) != 0 {
						macAddr := inter.HardwareAddr
						return ipnet.IP.String(), macAddr.String(), nil
					}
				}
			}
		}
	}

	resErr = errors.New(errString)
	return
}

func (p *ProAuthType) GetClientId() string {
	s, err := global.SystemSetting.GetValueString(ProAuthTypeSystemSettingClientIdKey)
	if err != nil {
		global.Logger.Debug("get clientId by syemtem setting:", err.Error())
		return ""
	}
	return s
}

func (p *ProAuthType) GetBase() proAuth.Base {
	ip, mac, err := p.GetLocalLanIpMAC()
	if err != nil {
		global.Logger.Error("get mac and ip error:", err.Error())
	}
	clientId := p.GetClientId()
	param := proAuth.Base{}
	param.Version = cmn.GetSysVersionInfo().Version
	param.MacAddress = mac
	param.LanIP = ip
	param.ClientId = &clientId
	param.Timestamp = time.Now().Unix()
	return param
}

func (p *ProAuthType) Register() {
	global.Logger.Debug("API: Register")
	param := proAuth.RegisterReq{}
	param.Base = p.GetBase()
	param.Code = uuid.NewString() + cmn.BuildRandCode(10, cmn.RAND_CODE_MODE1)

	resp, err := proAuth.Register(param)
	if err != nil {
		global.Logger.Error("register:", err.Error())
		return
	}

	// 设置clientId
	if resp.ClientId != "" {
		global.SystemSetting.Set(ProAuthTypeSystemSettingClientIdKey, resp.ClientId)
	}

	if param.ClientId == nil && resp.ClientId == "" {
		global.Logger.Error("register not fail: can not get clientId by server")
	}
}

func (p *ProAuthType) Login(username, password string) error {
	global.Logger.Debug("API: Login")
	param := proAuth.LoginReq{}
	param.Base = p.GetBase()
	param.Username = username
	param.Password = password
	p.LoginStatus = false

	if param.ClientId == nil || *param.ClientId == "" {
		p.Register()
	}

	resp, err := proAuth.Login(param)
	if err != nil {
		global.Logger.Error("login fail:", err.Error())
		return err
	}

	global.Logger.Debug(username+" token", resp.Token)

	if p.ProAuthDataShortTerm == nil {
		p.ProAuthDataShortTerm = &ProAuthDataShortTerm{}
	}
	p.ProAuthDataShortTerm.Token = resp.Token
	p.ProAuthDataShortTerm.ProExpiration = resp.ProExpiration
	p.ProAuthDataShortTerm.UserInfo = resp.UserInfo
	p.AuthDataShortTermSyncToSysntemSetting()

	p.AutoLoginDisable = false
	p.LoginStatus = true

	// 登录后 获取临时授权
	p.RenewTempAuth()

	return nil
}

func (p *ProAuthType) AutoLogin() error {
	global.Logger.Debug("API: AutoLogin")

	// 禁用自动登录
	if p.AutoLoginDisable {
		global.Logger.Debug("disable auto login")
		return errors.New("disable auto login")
	}

	accountPassword := p.GetAccountPassword()

	if accountPassword.Username == "" || accountPassword.Password == "" {
		p.AutoLoginDisable = true
		return ErrAccountPasswordIncomplete
	}

	param := proAuth.LoginReq{}
	param.Base = p.GetBase()
	param.Username = accountPassword.Username
	param.Password = accountPassword.Password

	if param.ClientId == nil || *param.ClientId == "" {
		p.Register()
	}

	resp, err := proAuth.AutoLogin(param)
	if err != nil {
		// p.LoginStatus = false
		if errors.Is(err, proAuth.GetErrorAndDefault(proAuth.ErrManualLogin, err)) {
			p.Reset()
		} else if errors.Is(err, proAuth.GetErrorAndDefault(proAuth.ErrIncorrectUsernameOrPassword, err)) {
			p.Reset()
		} else if errors.Is(err, proAuth.GetErrorAndDefault(proAuth.ErrAccountDisabled, err)) {
			p.Reset()
		}
		// p.AutoLoginDisable = true
		// 在授权期内，忽略其他错误
		global.Logger.Error("auto login fail:", err.Error())
		return err
	}

	global.Logger.Debug(param.Username+" token", resp.Token)

	if p.ProAuthDataShortTerm == nil {
		p.ProAuthDataShortTerm = &ProAuthDataShortTerm{}
	}
	p.ProAuthDataShortTerm.Token = resp.Token
	p.ProAuthDataShortTerm.ProExpiration = resp.ProExpiration
	p.ProAuthDataShortTerm.UserInfo = resp.UserInfo
	p.AuthDataShortTermSyncToSysntemSetting()

	p.AutoLoginDisable = false
	p.LoginStatus = true

	global.Logger.Debugln("自动邓罗", cmn.AnyToJsonStr(p))
	return nil
}

func (p *ProAuthType) RenewTempAuth() error {
	global.Logger.Debug("API: RenewTempAuth")
	param := proAuth.RenewTempAuthReq{}
	param.Base = p.GetBase()

	// 读取授权信息
	{
		shortAuth := p.GetAuthDataShortTerm()
		if shortAuth == nil {
			p.AutoLogin()
			return errors.New("need login")
		}
		param.Token = shortAuth.Token
	}

	if param.ClientId == nil || *param.ClientId == "" {
		p.Register()
	}

	resp, err := proAuth.RenewTempAuth(param)
	if err != nil {
		global.Logger.Error("renewTempAuth fail:", err.Error())

		// token过期,并且允许自动登录
		if errors.Is(err, proAuth.GetErrorAndDefault(proAuth.ErrTokenExpires, err)) && !p.AutoLoginDisable {
			p.AutoLogin()
			return nil
		}
		// 临时授权过期并连不上授权服务器，退出登录
		global.Logger.Info("LOGOUT:Short -term authorization expires and cannot be connected to the server, logout")
		p.Reset()
		return err
	}
	if p.ProAuthDataShortTerm == nil {
		p.ProAuthDataShortTerm = &ProAuthDataShortTerm{}
	}
	p.ProAuthDataShortTerm.ProExpiration = resp.ProExpiration
	p.ProAuthDataShortTerm.ShortTermExpiration = nil
	if resp.TempExpiration != "" {
		if t, err := time.Parse(cmn.TimeFormatMode1, resp.TempExpiration); err == nil {
			p.ProAuthDataShortTerm.ShortTermExpiration = &t
		}
	}

	p.AuthDataShortTermSyncToSysntemSetting()
	return nil
}

func (p *ProAuthType) Ping() error {
	global.Logger.Debug("API: Ping")
	param := proAuth.PingReq{}
	param.Base = p.GetBase()

	// 无需验证授权，正常是统计接口
	if p.ProAuthDataShortTerm == nil {
		p.ProAuthDataShortTerm = &ProAuthDataShortTerm{}
	}
	param.Token = p.ProAuthDataShortTerm.Token

	// ClientId不存在将注册设备
	if param.ClientId == nil || *param.ClientId == "" {
		p.Register()
	}

	param.Base = p.GetBase() // 重新获取ClientId
	resp, err := proAuth.Ping(param)
	if err != nil {
		if errors.Is(err, proAuth.GetErrorAndDefault(proAuth.ErrTokenExpires, err)) {
			global.Logger.Debug("ping fail:", err.Error())
			p.AutoLogin()
		}
		// 其他错误不处理
		return err
	}

	p.ProAuthDataShortTerm.ProExpiration = resp.ProExpiration

	return nil
}

func (p *ProAuthType) RefreshInfo() error {
	global.Logger.Debug("API: RefreshInfo")
	param := proAuth.RefreshInfoReq{}
	param.Base = p.GetBase()

	// 无需验证授权，正常是统计接口
	if p.ProAuthDataShortTerm == nil {
		p.ProAuthDataShortTerm = &ProAuthDataShortTerm{}
	}
	param.Token = p.ProAuthDataShortTerm.Token

	// ClientId不存在将注册设备
	if param.ClientId == nil || *param.ClientId == "" {
		p.Register()
	}

	respAny, err := p.requestAndAutoUpdateToken(func() (any, error) {
		param.Token = p.ProAuthDataShortTerm.Token
		// global.Logger.Debugln("refreshInfo JSON:", cmn.AnyToJsonStr(param))
		return proAuth.RefreshInfo(param)
	})

	if err != nil {
		// if errors.Is(err, proAuth.GetErrorAndDefault(proAuth.ErrManualLogin, err)) {
		// 	p.Reset()
		// } else if errors.Is(err, proAuth.GetErrorAndDefault(proAuth.ErrAccountDisabled, err)) {
		// 	p.Reset()
		// }
		// // p.AutoLoginDisable = true
		// // 在授权期内，忽略其他错误
		global.Logger.Error("refreshInfo fail-1:", err.Error())
		return err
	}

	if resp, ok := respAny.(proAuth.RefreshInfoResp); ok {
		p.ProAuthDataShortTerm.UserInfo = resp.UserInfo
		p.ProAuthDataShortTerm.ProExpiration = resp.ProExpiration
		return nil
	}

	respErr := errors.New("request data error")
	global.Logger.Error("refreshInfo fail-2:", respErr)

	return respErr
}

// 退出登录
func (p *ProAuthType) Reset() {
	p.SetAccountPassword("", "")
	p.AutoLoginDisable = true
	p.ProAuthDataShortTerm = nil
	p.LoginStatus = false
	p.AuthDataShortTermSyncToSysntemSetting()
}

func (p *ProAuthType) GetStatus() ProAuthType {
	return *p
}

func (p *ProAuthType) timer(d time.Duration, function func()) {
	// 创建一个定时器，每秒执行一次
	ticker := time.NewTicker(d)

	// 在main函数结束时停止定时器
	defer ticker.Stop()

	// 循环遍历定时器的通道
	for {
		select {
		case <-ticker.C:
			function()
		}
	}
}

func (p *ProAuthType) timerFunc() {
	p.Ping()

	// pro授权没有过期，但是临时授权过期了
	// 要进行续期临时授权

	proExpired := p.ProIsExpired()
	shortTermExpired := p.TempAuthIsExpired()

	if !proExpired && shortTermExpired {
		p.RenewTempAuth()
	} else if proExpired && shortTermExpired {
		global.Logger.Debug("PRO已过期")
	}
}

// 是否过期
func (p *ProAuthType) ProIsExpired() bool {
	global.Logger.Debug("检查是否过期 ProIsExpired")

	if p.ProAuthDataShortTerm == nil {
		return true
	}
	if p.ProAuthDataShortTerm.ProExpiration == nil {
		return true
	}
	global.Logger.Debug("临时授权过期时间", p.ProAuthDataShortTerm.ProExpiration)
	global.Logger.Debug("最终过期时间", p.ProAuthDataShortTerm.ProExpiration)
	return p.ProAuthDataShortTerm.ProExpiration.Before(time.Now())
}

func (p *ProAuthType) GetHideProBadgeStatus() bool {
	var status bool
	err := global.SystemSetting.GetValueByInterface("hide_pro_badge", &status)
	if err != nil {
		return false
	}

	return status
}

func (p *ProAuthType) SetHideProBadgeStatus(stats bool) {
	err := global.SystemSetting.Set("hide_pro_badge", stats)
	if err != nil {
		global.Logger.Errorln("SetHideProBadgeStatus", err)
	}
}

func (p *ProAuthType) TempAuthIsExpired() bool {
	global.Logger.Debug("检查是否过期 TempAuthIsExpired")
	if p.ProIsExpired() {
		return true
	}
	if p.ProAuthDataShortTerm == nil {
		return false
	}
	if p.ProAuthDataShortTerm.ShortTermExpiration == nil {
		return true
	}

	return p.ProAuthDataShortTerm.ShortTermExpiration.Before(time.Now())
}

func (p *ProAuthType) GetAccountPassword() ProAuthAccountPassword {
	accountPassword := ProAuthAccountPassword{}
	err := global.SystemSetting.GetValueByInterface(ProAuthTypeSystemSettingAccountPasswordKey, &accountPassword)
	if err != nil {
		global.Logger.Debug("GetAccountPassword error ", err.Error())
	}
	return accountPassword
}

func (p *ProAuthType) SetAccountPassword(username, password string) error {
	accountPassword := ProAuthAccountPassword{
		Username: username,
		Password: password,
	}
	return global.SystemSetting.Set(ProAuthTypeSystemSettingAccountPasswordKey, accountPassword)
}

// 保存临时（短期）授权数据
func (p *ProAuthType) SaveAuthDataShortTerm(proAuthDataShortTerm *ProAuthDataShortTerm) {
	d := ProAuthDataShortTerm{}
	if proAuthDataShortTerm != nil {
		d = *proAuthDataShortTerm
	}
	dataByte, _ := json.Marshal(d)
	dStr := string(dataByte)
	global.Logger.Debugln("AuthDataShortTerm JSON数据", dStr)

	// 明文
	// global.SystemSetting.Set(ProAuthTypeSystemSettingProAuthDataShortTermKey, dStr)

	// 加密
	aesStr, _ := AES.Encrypt(proAuth.Config.AESKey, dStr)
	global.SystemSetting.Set(ProAuthTypeSystemSettingProAuthDataShortTermKey, aesStr)

}

// 读取临时（短期）授权数据
func (p *ProAuthType) GetAuthDataShortTerm() *ProAuthDataShortTerm {
	d := ProAuthDataShortTerm{}
	dJsonStr := ""

	dEncryptedText, err := global.SystemSetting.GetValueString(ProAuthTypeSystemSettingProAuthDataShortTermKey)
	if err != nil {
		return nil
	}

	// 解密
	if dJsonStr, err = AES.Decrypt(proAuth.Config.AESKey, dEncryptedText); err != nil {
		return nil
	}

	// 明文
	// dJsonStr = dEncryptedText

	global.Logger.Debugln("AuthDataShortTerm 解密JSON数据", dJsonStr)

	err = json.Unmarshal([]byte(dJsonStr), &d)

	if d.Token == "" || err != nil {
		return nil
	}

	return &d

}

func (p *ProAuthType) AuthDataShortTermSyncToSysntemSetting() {
	p.SaveAuthDataShortTerm(p.ProAuthDataShortTerm)
}

func (p *ProAuthType) CheckNewVersion(versionType VersionType) (*LatestVersion, error) {
	var v *LatestVersion
	version, err := proAuth.CheckVersion(proAuth.CheckVersionReq{
		Base:        p.GetBase(),
		VersionType: proAuth.VersionType(versionType),
	})

	if err != nil {
		return nil, err
	}

	// jb, e := json.Marshal(version)
	// fmt.Println("检查更新返回的数据", string(jb), e)

	if version.IsNoNewVersion {
		return nil, nil
	}

	if version.VersionInfo != nil {
		v = &LatestVersion{
			ReleaseTime: version.VersionInfo.ReleaseTime,
			Type:        version.VersionInfo.Type,
			Version:     version.VersionInfo.Version,
			DownloadURL: version.VersionInfo.DownloadURL,
			PageUrl:     version.VersionInfo.PageUrl,
		}
	}

	return v, nil
}

func (p *ProAuthType) requestAndAutoUpdateToken(request func() (any, error)) (any, error) {
	// 第一次请求
	d, reqErr := request()
	if reqErr == nil {
		return d, nil // 请求成功，返回数据和 nil 错误
	}

	// 如果错误是 token 过期
	if errors.Is(reqErr, proAuth.GetErrorAndDefault(proAuth.ErrTokenExpires, reqErr)) {
		// 自动更新 token
		autoErr := p.AutoLogin()
		if autoErr != nil {
			return nil, autoErr // 自动登录失败，返回错误
		}

		global.Logger.Errorln("自动重新登录成功，准备二次调用")
		// 自动登录成功后，重新尝试请求
		d, reqErr = request()
		if reqErr != nil {
			return nil, reqErr // 第二次请求仍然失败，返回错误
		}
		return d, nil // 成功返回数据和 nil 错误
	}

	// 其他错误直接返回
	return nil, reqErr
}
