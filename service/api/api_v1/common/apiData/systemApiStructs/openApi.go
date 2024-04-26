package systemApiStructs

type OpenApiTokenInfo struct {
	Token string `json:"token"`
}

type OpenApiGetTokenStatusInfoResp struct {
	OpenApiTokenInfo
	Status bool `json:"status"`
}
