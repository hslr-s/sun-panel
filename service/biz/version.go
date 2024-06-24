package biz

import (
	"sun-panel/global"
	"time"
)

type VersionStruct struct {
}

var (
	VersionAutoCheckSystemConfigName       = "versionAutoCheckSystemConfig" // 自动检查更新的配置名称
	VersionNewVersionCacheSystemConfigName = "versionNewVersionCache"       // 新版本缓存信息配置名称
)

// 记录是否自动更新的配置
type VersionAutoSystemConfig struct {
	Beta    bool `json:"beta"`
	Release bool `json:"release"`
}

type VersionNewVersionCacheSystemConfigIgnored struct {
	Version    string
	UpdateTime time.Time
}

// 记录检查新版本缓存的信息
type VersionNewVersionCacheSystemConfig struct {
	VersionNewVersionCacheSystemConfigIgnored []VersionNewVersionCacheSystemConfigIgnored
	BetaLatest                                LatestVersion
	ReleaseLatest                             LatestVersion
}

func (p *VersionStruct) SetAutoCheckNewVersionConfig(cfg VersionAutoSystemConfig) error {
	return global.SystemSetting.Set(VersionAutoCheckSystemConfigName, cfg)
}

func (p *VersionStruct) GetAutoCheckNewVersionConfig() (VersionAutoSystemConfig, error) {
	v := VersionAutoSystemConfig{}
	err := global.SystemSetting.GetValueByInterface(VersionAutoCheckSystemConfigName, &v)
	return v, err
}

func (p *VersionStruct) SetVersionNewVersionCacheSystemConfig(cfg VersionAutoSystemConfig) error {
	return global.SystemSetting.Set(VersionNewVersionCacheSystemConfigName, cfg)
}

func (p *VersionStruct) GetVersionNewVersionCacheSystemConfig() (VersionAutoSystemConfig, error) {
	v := VersionAutoSystemConfig{}
	err := global.SystemSetting.GetValueByInterface(VersionNewVersionCacheSystemConfigName, &v)
	return v, err
}
