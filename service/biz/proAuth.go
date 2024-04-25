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

var ErrAccountPasswordIncomplete = errors.New("ErrAccountPasswordIncomplete")

const ProAuthTypeSystemSettingClientIdKey = "clientId"

// 短期授权的key
const ProAuthTypeSystemSettingProAuthDataShortTermKey = "pro_auth_data_short_term"

// 将废弃 (1.4.0-beta24-04-12添加的此项)
const ProAuthTypeSystemSettingAccountPasswordKey = "ProAuthAccountPassword"

func NewProAuth() *ProAuthType {

	proAuth.Config.Url = "https://auth.sun-panel.top"
	proAuth.Config.AESKey = "wAOhhAqbBxgjNOlcEAGvwawvUZmWDLkN"
	proAuth.Config.TimeDifference = 3600

	p := ProAuthType{
		PingInterval:     6 * time.Hour, // ping执行的时间间隔
		AutoLoginDisable: false,         // 刚启动允许自动登录
	}

	if global.RUNCODE == "debug" {
		proAuth.Config.Url = "http://127.0.0.1:3005"
		p.PingInterval = 30 * time.Second
		// proAuth.Config.IsDebug = true // 明文传输
	}

	return &p
}

func (p *ProAuthType) Start() {
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
	p.SaveAuthDataShortTerm(*p.ProAuthDataShortTerm)

	p.AutoLoginDisable = false
	return nil
}

func (p *ProAuthType) AutoLogin() error {
	global.Logger.Debug("API: AutoLogin")

	// 禁用自动登录
	if p.AutoLoginDisable {
		global.Logger.Info("disable auto login")
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
		if errors.Is(err, proAuth.GetError(proAuth.ErrManualLogin)) {
			p.AutoLoginDisable = true
		} else if errors.Is(err, proAuth.GetError(proAuth.ErrAccountNotExist)) {
			p.AutoLoginDisable = true
		} else if errors.Is(err, proAuth.GetError(proAuth.ErrIncorrectUsernameOrPassword)) {
			p.AutoLoginDisable = true
		} else if errors.Is(err, proAuth.GetError(proAuth.ErrAccountDisabled)) {
			p.AutoLoginDisable = true
		}

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
	p.SaveAuthDataShortTerm(*p.ProAuthDataShortTerm)

	p.AutoLoginDisable = false
	return nil
}

func (p *ProAuthType) RenewTempAuth() error {
	global.Logger.Debug("API: RenewTempAuth")
	param := proAuth.RenewTempAuthReq{}
	param.Base = p.GetBase()

	// 读取授权信息
	{
		shortAuth, err := p.GetAuthDataShortTerm()
		if err != nil {
			return err
		}
		param.Token = shortAuth.Token
	}

	if param.ClientId == nil || *param.ClientId == "" {
		p.Register()
	}

	resp, err := proAuth.RenewTempAuth(param)
	if err != nil {
		global.Logger.Error("renewTempAuth fail:", err.Error())
		return err
	}
	p.ProExpiration = resp.ProExpiration
	p.TempExpiration = resp.TempExpiration

	// 保存授权信息
	shortAuth := ProAuthDataShortTerm{
		Token:               resp.Token,
		UserInfo:            p.UserInfo,
		ProExpiration:       resp.ProExpiration,
		ShortTermExpiration: resp.TempExpiration,
	}
	p.SaveAuthDataShortTerm(shortAuth)
	return nil
}

func (p *ProAuthType) Ping() error {
	global.Logger.Debug("API: Ping")
	param := proAuth.PingReq{}
	param.Base = p.GetBase()
	param.Token = p.Token

	// ClientId不存在将注册
	if param.ClientId == nil || *param.ClientId == "" {
		p.Register()
	}
	resp, err := proAuth.Ping(param)
	if err != nil {
		if errors.Is(err, proAuth.GetError(proAuth.ErrTokenExpires)) {
			global.Logger.Debug("ping fail:", err.Error())
			p.AutoLogin()
		}
		return err
	}

	p.ProExpiration = resp.ProExpiration

	return nil
}

func (p *ProAuthType) Reset() {
	p.SetAccountPassword("", "")
	p.AutoLoginDisable = true
	p.Token = ""
	p.TempExpiration = ""
	p.ProExpiration = nil
	p.UserInfo = proAuth.UserInfo{
		Username: "",
		Name:     "",
	}
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
	if !p.ProIsExpired() && p.TempAuthIsExpired() {
		p.RenewTempAuth()
	}
}

// 是否过期
func (p *ProAuthType) ProIsExpired() bool {
	if p.ProExpiration == nil {
		return true
	}
	return p.ProExpiration.Before(time.Now())
}

func (p *ProAuthType) TempAuthIsExpired() bool {
	if p.TempExpiration == "" {
		return true
	}
	timeParse, err := time.Parse(cmn.TimeFormatMode1, p.TempExpiration)
	if err != nil {
		global.Logger.Error("TempAuthIsExpired:TempExpiration parse error", timeParse)
		return true
	}
	return timeParse.Before(time.Now())
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
func (p *ProAuthType) SaveAuthDataShortTerm(proAuthDataShortTerm ProAuthDataShortTerm) {
	dataByte, _ := json.Marshal(proAuthDataShortTerm)
	dStr := string(dataByte)
	global.Logger.Debugln("AuthDataShortTerm JSON数据", dStr)
	aesStr, _ := AES.Encrypt(proAuth.Config.AESKey, dStr)
	global.SystemSetting.Set(ProAuthTypeSystemSettingAccountPasswordKey, aesStr)
}

// 读取临时（短期）授权数据
func (p *ProAuthType) GetAuthDataShortTerm() (ProAuthDataShortTerm, error) {
	d := ProAuthDataShortTerm{}
	dJsonStr := ""

	dEncryptedText, err := global.SystemSetting.GetValueString(ProAuthTypeSystemSettingAccountPasswordKey)
	if err != nil {
		return d, err
	}

	if dJsonStr, err = AES.Decrypt(proAuth.Config.AESKey, dEncryptedText); err != nil {
		return d, err
	}

	global.Logger.Debugln("AuthDataShortTerm 解密JSON数据", dJsonStr)

	json.Unmarshal([]byte(dJsonStr), &d)

	return d, err

}
