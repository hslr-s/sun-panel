package biz

import (
	"errors"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/models"

	"gorm.io/gorm"
)

type OpanApiType struct {
}

func (o *OpanApiType) GetUserInfo(token string) (openTokenInfo models.OpenApiToken, gormErr error) {
	err := global.Db.Preload("User").Order("id DESC, created_at DESC").First(&openTokenInfo, "token=?", token).Error
	return openTokenInfo, err
}

func (o *OpanApiType) GetUserTokenByUserId(userId uint) (openTokenInfo models.OpenApiToken, gormErr error) {
	err := global.Db.Order("id DESC, created_at DESC").First(&openTokenInfo, "user_id=?", userId).Error
	return openTokenInfo, err
}

func (o *OpanApiType) RefreshToken(userId uint) (openTokenInfo string, gormErr error) {
	newToken := ""
	buildTokenOver := false
	for !buildTokenOver {
		newToken = o.buildToken()
		global.Logger.Debugln("newToken", newToken)
		if _, err := o.GetUserInfo(newToken); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 创建token
				errCreate := global.Db.Transaction(func(tx *gorm.DB) error {
					// 删除现有的token
					tx.Delete(&models.OpenApiToken{}, "user_id=?", userId)

					// 创建一个新的token
					err := tx.Create(&models.OpenApiToken{
						UserId: userId,
						Token:  newToken,
					}).Error

					if err != nil {
						return err
					}
					return nil
				})
				return newToken, errCreate
			}
			buildTokenOver = true
		}

	}
	return newToken, nil
}

func (o *OpanApiType) DeleteToken(userId uint) error {
	return global.Db.Delete(&models.OpenApiToken{}, "user_id=?", userId).Error
}

func (o *OpanApiType) buildToken() string {
	return cmn.BuildRandCode(32, cmn.RAND_CODE_MODE2)
}
