package config

import (
	"os"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/lib/iniConfig"
)

func getDefaultConfig() map[string]map[string]string {
	return map[string]map[string]string{
		"base": {
			"http_port":        "9090",
			"source_path":      "./files",      // 存放文件的路径
			"source_temp_path": "./files/temp", // 存放文件的缓存路径
		},
		"sqlite": {
			"file_path": "./database.db",
		},
	}

}

func ConfigInit() (*iniConfig.IniConfig, error) {

	// 配置文件初始化
	if config, err, errCode := Conf(getDefaultConfig()); err != nil && errCode == 0 {
		// 抛出错误
		cmn.Pln(cmn.LOG_ERROR, "配置文件创建错误:"+err.Error())
		os.Exit(1)
		return nil, err
	} else if errCode == 1 {
		// 配置文件不存在，进行创建
		if err := CreateConfExample("conf.example.ini", "conf.ini"); err != nil {
			cmn.Pln(cmn.LOG_ERROR, "配置文件创建错误:"+err.Error())
			os.Exit(1)
			return nil, err
		}

		global.Logger.Infoln("配置文件已经自动生成'conf/conf.ini',将再次读取配置")
		// 创建成功再次读取文件
		if configAgain, errAgain, _ := Conf(getDefaultConfig()); errAgain != nil {
			return nil, errAgain
		} else {
			global.Logger.Infoln("尝试读取配置文件'conf/conf.ini',二次读取配置文件成功")
			return configAgain, nil
		}
	} else {
		return config, nil
	}
}

// 配置初始化
// errCode=1 说明初始化流程
func Conf(defaultConfig map[string]map[string]string) (config *iniConfig.IniConfig, err error, errCode int) {
	CreateConfExample("conf.example.ini", "conf.example.ini")
	exists, err := cmn.PathExists("conf/conf.ini")
	if exists {
		config = iniConfig.NewIniConfig("conf/conf.ini") // 读取配置
		config.Default = defaultConfig
	} else if err != nil {

	} else {
		// docker 运行模式，生成配置文件
		if global.ISDOCKER != "" {
			cmn.AssetsTakeFileToPath("conf.example.ini", "conf/conf.ini")
			config = iniConfig.NewIniConfig("conf/conf.ini") // 读取配置
			config.Default = defaultConfig
		} else {
			errCode = 1
		}
	}
	return
}

// 生成示例配置文件
func CreateConfExample(confName string, targetName string) (err error) {
	// 查看配置示例文件是否存在，不存在创建（分别为示例配置和配置文件）
	exists, err := cmn.PathExists("conf/" + targetName)
	if err != nil {
		return
	}
	if !exists {
		if err = cmn.AssetsTakeFileToPath(confName, "conf/"+targetName); err != nil {
			return
		}
	}

	return nil
}
