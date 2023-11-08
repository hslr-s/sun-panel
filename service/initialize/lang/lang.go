package lang

import (
	"os"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/lib/language"
)

func LangInit(lang string) {
	filename := "lang/" + lang + ".ini"
	exists, err := cmn.PathExists(filename)
	if err != nil {
		global.Logger.Errorln("语言文件不存在", err.Error())
		os.Exit(1)
	}

	// 生成语言文件
	if !exists {
		global.Logger.Infoln("输出语言文件:", filename)
		err := cmn.AssetsTakeFileToPath("lang/zh-cn.ini", "lang/zh-cn.ini")
		if err != nil {
			global.Logger.Errorln("输出语言文件出错:", err.Error())
			os.Exit(1)
		}
		err = cmn.AssetsTakeFileToPath("lang/en-us.ini", "lang/en-us.ini")
		if err != nil {
			global.Logger.Errorln("输出语言文件出错:", err.Error())
			os.Exit(1)
		}
	}
	exists, err = cmn.PathExists(filename)
	if err != nil || !exists {
		global.Logger.Errorln("语言文件不存在:", filename)
		os.Exit(1)
	}

	global.Lang = language.NewLang(filename)
}
