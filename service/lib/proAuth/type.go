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
type RegisterReq struct {
	Base
	Code string // 随机码，避免重复攻击
}

type RespBase struct {
	Token     string
	Timestamp int64
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
