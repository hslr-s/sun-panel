package config

import (
	"bytes"
	"os"
	"path"
	"sun-panel/assets"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/lib/iniConfig"
)

var (
	GlobalConfigPath = "conf"
)

// 初始化配置文件, 并返回配置文件对象
// 正常情况`conf.example.ini` 是默认覆盖原有的 `conf.example.ini`。
// `conf.ini` 则是根据 'coverageConfIni' 参数来决定是否覆盖
func ConfigInit(cfgPath string, coverageConfIni bool) (*iniConfig.IniConfig, error) {

	// 创建示例配置文件,强制覆盖最新的示例配置文件
	if err := CreateConf("conf.example.ini", cfgPath+"/conf.example.ini", true); err != nil {
		// 示例配置文件生成失败
		return nil, err
	}

	// 生成配置文件,如果文件存在忽略
	if err := CreateConf("conf.example.ini", cfgPath+"/conf.ini", coverageConfIni); err != nil {
		return nil, err
	}

	// 读取配置文件
	config := iniConfig.NewIniConfig(cfgPath + "/conf.ini")
	return config, nil
}

// 生成配置文件
func CreateConf(assetsConfName string, targetFile string, forced bool) error {
	// 如果配置文件存在,选择是否强制覆盖
	exists, err := cmn.PathExists(targetFile)
	if err != nil {
		return err
	}

	if !exists || (exists && forced) {
		// 生成配置文件

		// 创建目录
		if err := os.MkdirAll(path.Dir(targetFile), 0777); err != nil {
			return err
		}

		if originalBytes, err := assets.Asset("assets/" + assetsConfName); err != nil {
			return err
		} else {
			// 替换配置文件中的变量 __$PATH__
			oldStr := []byte("__$PATH__")
			newStr := []byte(".")

			// 为了减少挂载docker就只存放在conf文件夹
			if global.ISDOCKER == "true" {
				newStr = []byte("./conf")
			}

			global.Logger.Debugln("当前环境是否为docker:", global.ISDOCKER)
			configFileBytes := bytes.ReplaceAll(originalBytes, oldStr, newStr)
			return os.WriteFile(targetFile, configFileBytes, 0666)
		}
	}

	// 不需要创建
	return nil
}
