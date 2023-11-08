package adminApiStructs

type GetStatisticsResp struct {
	UserCount   int64 `json:"userCount"`
	UserToday   int64 `json:"userToday"`
	RoleCount   int64 `json:"roleCount"`
	RoleToday   int64 `json:"roleToday"`
	DialogCount int64 `json:"dialogCount"`
	DialogToday int64 `json:"dialogToday"`
	DrawCount   int64 `json:"drawCount"`
	DrawToday   int64 `json:"drawToday"`
}
