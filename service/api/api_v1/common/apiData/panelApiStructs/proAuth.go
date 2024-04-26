package panelApiStructs

import "time"

type LoginReq struct {
	Username string
	Password string
}

type GetStatusResp struct {
	ProExpiration *time.Time `json:"proExpiration"`
	LoginStatus   bool       `json:"loginStatus"`
	Username      string     `json:"username"`
	NickName      string     `json:"nickName"`
}
