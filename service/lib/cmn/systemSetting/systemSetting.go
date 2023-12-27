package systemSetting

import (
	"encoding/json"
	"errors"
	"sun-panel/lib/cache"
	"sun-panel/models"

	"gorm.io/gorm"
)

const (
	SYSTEM_APPLICATION    = "system_application"
	SYSTEM_EMAIL          = "system_email"
	DISCLAIMER            = "disclaimer"            // 免责声明 储存类型：字符串
	WEB_ABOUT_DESCRIPTION = "web_about_description" // 关于的描述信息
	PANEL_PUBLIC_USER_ID  = "panel_public_user_id"  // 公开访问模式用户id *uint|null
)

type SystemSettingCache struct {
	Cache cache.Cacher[interface{}]
}

type Email struct {
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	Mail     string `json:"mail" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type Register struct {
	EmailSuffix  string `json:"emailSuffix"`  // 注册邮箱后缀
	OpenRegister bool   `json:"openRegister"` // 开放注册
}

type Login struct {
	LoginCaptcha bool `json:"loginCaptcha"` // 登录验证码
}

type ApplicationSetting struct {
	Register
	Login
	WebSiteUrl string `json:"webSiteUrl"` // 站点地址
}

var (
	ErrorNoExists = errors.New("no exists")
)

// 系统配置启用缓存功能
func (s *SystemSettingCache) GetValueString(configName string) (result string, err error) {
	if v, ok := s.Cache.Get(configName); ok {
		if v1, ok1 := v.(string); ok1 {
			// fmt.Println("读取缓存")
			return v1, nil
		}
	}

	mSetting := models.SystemSetting{}
	result, err = mSetting.Get(configName)
	if err == gorm.ErrRecordNotFound {
		err = ErrorNoExists
	}
	// 查询出来，缓存起来
	s.Cache.SetDefault(configName, result)
	return
}

// value 需为指针
func (s *SystemSettingCache) GetValueByInterface(configName string, value interface{}) error {
	if v, ok := s.Cache.Get(configName); ok {
		// fmt.Println("缓存")
		if s, sok := v.(string); sok {
			if err := json.Unmarshal([]byte(s), value); err != nil {
				return err
			}
			return nil
		}
	}

	mSetting := models.SystemSetting{}
	result, err := mSetting.Get(configName)
	if err == gorm.ErrRecordNotFound {
		err = ErrorNoExists
		return err
	}
	err = json.Unmarshal([]byte(result), value)
	if err != nil {
		return err
	}
	s.Cache.SetDefault(configName, result)
	return nil
}

func (s *SystemSettingCache) Set(configName string, configValue interface{}) error {
	s.Cache.Delete(configName)
	mSetting := models.SystemSetting{}
	err := mSetting.Set(configName, configValue)
	return err
}
