package language

import (
	"os"
	"strings"
	"sun-panel/lib/cmn"
	"sun-panel/lib/iniConfig"
)

type LangStructObj struct {
	LangContet *iniConfig.IniConfig
}

func NewLang(langPath string) *LangStructObj {
	langObj := LangStructObj{}
	exists, _ := cmn.PathExists(langPath)

	if exists {
		langObj.LangContet = iniConfig.NewIniConfig(langPath) // 读取配置
	} else {
		cmn.Pln(cmn.LOG_ERROR, "language file does not exist:"+langPath)
		os.Exit(1)
	}
	return &langObj
}

// 获取
// common.lang
// 配置文件格式
// [common]
// lang=zh-cn
func (l *LangStructObj) Get(key string) string {
	if key == "" {
		return key
	}
	keyArr := strings.Split(key, ".")
	if len(keyArr) < 2 {
		return l.LangContet.GetValueString(keyArr[0], "NOT EMPTY")
	} else {
		return l.LangContet.GetValueString(keyArr[0], keyArr[1])
	}
}

// 获取并替换字段
func (l *LangStructObj) GetWithFields(key string, fields map[string]string) string {
	c := l.Get(key)
	for k, v := range fields {
		c = strings.ReplaceAll(c, `{`+k+`}`, v)
	}
	return c
}

// 获取值并向后追加
func (l *LangStructObj) GetAndInsert(key string, insertContent ...string) string {
	content := l.Get(key) + " "
	for _, v := range insertContent {
		content += v
	}
	return content
}
