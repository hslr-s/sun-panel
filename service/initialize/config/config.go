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

func ConfigInit() (*iniConfig.IniConfig, error) {

	confPath := "./conf"

	// 创建示例配置文件,强制覆盖最新的示例配置文件
	if err := CreateConf("conf.example.ini", confPath+"/conf.example.ini", true); err != nil {
		// 示例配置文件生成失败
		return nil, err
	}

	// 生成配置文件,如果文件存在忽略
	if err := CreateConf("conf.example.ini", confPath+"/conf.ini", false); err != nil {
		return nil, err
	}

	// 读取配置文件
	config := iniConfig.NewIniConfig(confPath + "/conf.ini")
	return config, nil

	// // 配置文件初始化
	// if config, err, errCode := Conf(getDefaultConfig()); err != nil && errCode == 0 {
	// 	// 抛出错误
	// 	cmn.Pln(cmn.LOG_ERROR, "配置文件创建错误:"+err.Error())
	// 	os.Exit(1)
	// 	return nil, err
	// } else if errCode == 1 {
	// 	// 配置文件不存在，进行创建
	// 	if err := CreateConf("conf.example.ini", "conf.ini"); err != nil {
	// 		cmn.Pln(cmn.LOG_ERROR, "配置文件创建错误:"+err.Error())
	// 		os.Exit(1)
	// 		return nil, err
	// 	}

	// 	global.Logger.Infoln("配置文件已经自动生成'conf/conf.ini',将再次读取配置")
	// 	// 创建成功再次读取文件
	// 	if configAgain, errAgain, _ := Conf(getDefaultConfig()); errAgain != nil {
	// 		return nil, errAgain
	// 	} else {
	// 		global.Logger.Infoln("尝试读取配置文件'conf/conf.ini',二次读取配置文件成功")
	// 		return configAgain, nil
	// 	}
	// } else {
	// 	return config, nil
	// }
}

// 配置初始化
// errCode=1 说明初始化流程
// func Conf(defaultConfig map[string]map[string]string) (config *iniConfig.IniConfig, err error, errCode int) {
// 	CreateConf("conf.example.ini", "conf.example.ini")
// 	exists, err := cmn.PathExists("conf/conf.ini")
// 	if exists {
// 		config = iniConfig.NewIniConfig("conf/conf.ini") // 读取配置
// 		config.Default = defaultConfig
// 	} else if err != nil {

// 	} else {
// 		// docker 运行模式，生成配置文件
// 		if global.ISDOCKER != "" {
// 			cmn.AssetsTakeFileToPath("conf.example.ini", "conf/conf.ini")
// 			config = iniConfig.NewIniConfig("conf/conf.ini") // 读取配置
// 			config.Default = defaultConfig
// 		} else {
// 			errCode = 1
// 		}
// 	}
// 	return
// }

// 生成示例配置文件
// func CreateConfExample(confName string, targetName string) (err error) {
// 	// 查看配置示例文件是否存在，不存在创建（分别为示例配置和配置文件）
// 	exists, err := cmn.PathExists("conf/" + targetName)
// 	if err != nil {
// 		return
// 	}
// 	if !exists {
// 		if err = cmn.AssetsTakeFileToPath(confName, "conf/"+targetName); err != nil {
// 			return
// 		}
// 	}

// 	return nil
// }

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
			configFileBytes := bytes.Replace(originalBytes, oldStr, newStr, -1)
			return os.WriteFile(targetFile, configFileBytes, 0666)
		}
	}

	// 不需要创建
	return nil
}
