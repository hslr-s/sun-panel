package proAuth

import "time"

type Base struct {
	ClientId   *string
	Version    string
	MacAddress string
	LanIP      string
	Timestamp  int64
}

type TokenBase struct {
	Base
	Token string
}

type RespBase struct {
	Token     string
	Timestamp int64
}

type UserInfo struct {
	Username string
	Name     string
}

type RegisterReq struct {
	Base
	Code string // 随机码，避免重复攻击
}
type RegisterResp struct {
	RespBase
	ClientId string
}

type LoginReq struct {
	Base
	Username string
	Password string
}

type LoginResp struct {
	RespBase
	UserInfo
	Token         string
	ProExpiration *time.Time // pro授权过期时间
}

type RenewTempAuthReq struct {
	TokenBase
}

type RenewTempAuthResp struct {
	RespBase
	ProExpiration  *time.Time // pro授权过期时间
	TempExpiration string     // 临时授权过期时间（客户端超过这个过期时间将续期短期授权）
}

type PingReq struct {
	TokenBase
}

type PingResp struct {
	RespBase
	ProExpiration *time.Time // pro授权过期时间
}

type RefreshInfoReq struct {
	TokenBase
}

type RefreshInfoResp struct {
	UserInfo
	RespBase
	ProExpiration *time.Time // pro授权过期时间
}

type VersionType string

const (
	VersionTypeRelease VersionType = "release"
	VersionTypeBeta    VersionType = "beta"
	VersionTypeAlpha   VersionType = "alpha"
	VersionTypeRC      VersionType = "rc"
)

type CheckVersionReq struct {
	Base
	VersionType VersionType // 查询版本的类型
}
type VersionInfo struct {
	Version      string
	Type         VersionType
	ReleaseTime  time.Time
	Description  string
	DownloadURL  string
	PageUrl      string
	IsActive     bool
	IsRolledBack bool
}
type CheckVersionResp struct {
	RespBase
	VersionInfo    *VersionInfo
	IsNoNewVersion bool
}
