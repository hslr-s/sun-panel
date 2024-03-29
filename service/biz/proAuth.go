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
	ProExpiration  *time.Time
	TempExpiration string
	Token          string
}

const ProAuthTypeSystemSettingKey = "clientId"

func NewProAuth() *ProAuthType {
	proAuth.Config.Url = "http://127.0.0.1:3005"
	// proAuth.Config.Url = "http://127.0.0.1:8888"
	proAuth.Config.AESKey = "wAOhhAqbBxgjNOlcEAGvwawvUZmWDLkN"
	proAuth.Config.IsDebug = true
	return &ProAuthType{}
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
	s, err := global.SystemSetting.GetValueString(ProAuthTypeSystemSettingKey)
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
		global.SystemSetting.Set(ProAuthTypeSystemSettingKey, resp.ClientId)
	}

	if param.ClientId == nil && resp.ClientId == "" {
		global.Logger.Error("register not fail: can not get clientId by server")
	}
}

func (p *ProAuthType) Login(username, password string) error {
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

	global.Logger.Debug("token", resp.Token)

	p.Token = resp.Token
	p.ProExpiration = resp.ProExpiration
	return nil
}

func (p *ProAuthType) RenewTempAuth() error {
	param := proAuth.RenewTempAuthReq{}
	param.Base = p.GetBase()
	param.Token = &p.Token

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
	param := proAuth.PingReq{}
	param.Base = p.GetBase()
	param.Token = &p.Token

	if param.ClientId == nil || *param.ClientId == "" {
		p.Register()
	}

	resp, err := proAuth.Ping(param)
	if err != nil {
		global.Logger.Error("ping fail:", err.Error())
		return err
	}

	p.ProExpiration = resp.ProExpiration

	return nil
}
