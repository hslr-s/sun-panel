package biz

import (
	"sun-panel/global"
	"sun-panel/models"
)

type OpanApiType struct {
}

func (o *OpanApiType) GetUserInfo(token string) (openTokenInfo models.OpenApiToken, gormErr error) {
	err := global.Db.Preload("User").First(&openTokenInfo, "token=?", token).Error
	return openTokenInfo, err
}
