package biz

import (
	"errors"
	"net"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/lib/proAuth"
	"time"

	"github.com/google/uuid"
)

type ProAuthType struct {
	ProExpiration    *time.Time
	TempExpiration   string
	Token            string
	AutoLoginDisable bool // 为true时需要手动登录
	PingInterval     time.Duration
}

type ProAuthAccountPassword struct {
	Username string
	Password string
}

var ErrAccountPasswordIncomplete = errors.New("ErrAccountPasswordIncomplete")

const ProAuthTypeSystemSettingClientIdKey = "clientId"
const ProAuthTypeSystemSettingAccountPasswordKey = "ProAuthAccountPassword"

func NewProAuth() *ProAuthType {

	proAuth.Config.Url = "http://127.0.0.1:3005"
	// proAuth.Config.Url = "http://127.0.0.1:8888"
	proAuth.Config.AESKey = "wAOhhAqbBxgjNOlcEAGvwawvUZmWDLkN"
	// proAuth.Config.IsDebug = true

	p := ProAuthType{
		PingInterval:     30 * time.Second, // ping执行的时间间隔
		AutoLoginDisable: false,            // 刚启动允许自动登录
	}
	return &p
}

func (p *ProAuthType) Start() {
	// p.Login("95302870@qq.com", "12345678")
	p.Ping()
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
		global.Logger.Error("get clientId by syemtem setting:", err.Error())
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

	p.Token = resp.Token
	p.ProExpiration = resp.ProExpiration
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
		} else if errors.Is(err, proAuth.GetError(proAuth.ErrAccountNotExist)) {
			p.AutoLoginDisable = true
		} else if errors.Is(err, proAuth.GetError(proAuth.ErrAccountDisabled)) {
			p.AutoLoginDisable = true
		}

		global.Logger.Error("auto login fail:", err.Error())
		return err
	}

	global.Logger.Debug(param.Username+" token", resp.Token)

	p.Token = resp.Token
	p.ProExpiration = resp.ProExpiration
	return nil
}

func (p *ProAuthType) RenewTempAuth() error {
	global.Logger.Debug("API: RenewTempAuth")
	param := proAuth.RenewTempAuthReq{}
	param.Base = p.GetBase()
	param.Token = p.Token

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

func (p *ProAuthType) UpdateLoginStatus(username, password string) {

}
