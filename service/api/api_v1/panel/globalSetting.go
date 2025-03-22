package panel

import (
	"os"
	"sun-panel/api/api_v1/common/apiData/panelApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/lib/cmn"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type GlobalSettingApi struct {
}

func (a *GlobalSettingApi) getCustomPath() string {
	webPath := global.Config.GetValueStringOrDefault("base", "web_path")
	customPath := global.Config.GetValueString("base", "custom_style_path")
	if customPath == "" {
		customPath = webPath + "/custom"
	}

	return customPath
}

// 保存全局站点设置
func (a *GlobalSettingApi) SaveSiteStting(c *gin.Context) {
	req := panelApiStructs.GlobalSettingSite{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	biz.GlobalSetting.SaveSiteStting(req.GlobalSettingSiteSetting)
	apiReturn.Success(c)
}

func (a *GlobalSettingApi) GetSiteStting(c *gin.Context) {
	resp := panelApiStructs.GlobalSettingSite{}

	biz.GlobalSetting.GetSiteSttingNoCheckExpired(&resp)
	apiReturn.SuccessData(c, resp)
}

// 获取全局站点设置
func (a *GlobalSettingApi) UploadFaviconImage(c *gin.Context) {
	fileInfo, err := base.UploadFileAndSave(c, []string{
		".ico",
	}, 0)

	if err != nil {
		switch err {
		case base.ErrUploadExtensionNameNotAllowed:
			apiReturn.ErrorByCode(c, 1301)
		default:
			global.Logger.Error(err.Error())
			apiReturn.ErrorByCode(c, 1300)
		}
		return
	}

	apiReturn.SuccessData(c, gin.H{
		"imageUrl": global.UPLOAD_ROUTE + fileInfo.Src,
	})
}

func (a *GlobalSettingApi) GetCustomStylePath(c *gin.Context) {
	customPath := a.getCustomPath()
	apiReturn.SuccessData(c, gin.H{
		"customPath": customPath,
	})
}

func (a *GlobalSettingApi) GetCustomJsAndCssCode(c *gin.Context) {
	customPath := a.getCustomPath()
	resp := panelApiStructs.GlobalSettingCustomJsAndCss{}

	cssFilePath := customPath + "/index.css"
	cssIsExist, _ := cmn.PathExists(customPath)
	if cssIsExist {
		if s, err := readFile(cssFilePath); err == nil {
			resp.CssContent = s
		}
	}

	jsFilePath := customPath + "/index.js"
	jsIsExist, _ := cmn.PathExists(customPath)
	if jsIsExist {
		if s, err := readFile(jsFilePath); err == nil {
			resp.JsContent = s
		}
	}
	apiReturn.SuccessData(c, resp)
}

func (a *GlobalSettingApi) SaveCustomJsAndCssCode(c *gin.Context) {
	customPath := a.getCustomPath()

	if ok, err := cmn.PathExists(customPath); err != nil || !ok {
		err := os.MkdirAll(customPath, os.ModePerm)
		if err != nil {
			apiReturn.ErrorByCode(c, apiReturn.ErrCodeCustomJSAndCSSFilePreservationFailure)
			global.Logger.Errorln("failure to create a custom directory failed:", customPath, " - ", err.Error())
		}
	}

	req := panelApiStructs.GlobalSettingCustomJsAndCss{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if biz.ProAuth.ProIsExpired() {
		apiReturn.ErrorByCode(c, apiReturn.ErrCodeNoPROAuthorization)
		return
	}

	cssFilePath := customPath + "/index.css"
	// cssIsExist, _ := cmn.PathExists(customPath)
	if err := createAndWriteFile(cssFilePath, req.CssContent); err != nil {
		apiReturn.ErrorByCode(c, apiReturn.ErrCodeCustomJSAndCSSFilePreservationFailure)
		global.Logger.Errorln("createAndWriteFile:", cssFilePath, " - ", err.Error())
		return
	}

	jsFilePath := customPath + "/index.js"
	// cssIsExist, _ := cmn.PathExists(customPath)
	if err := createAndWriteFile(jsFilePath, req.JsContent); err != nil {
		apiReturn.ErrorByCode(c, apiReturn.ErrCodeCustomJSAndCSSFilePreservationFailure)
		global.Logger.Errorln("createAndWriteFile:", cssFilePath, " - ", err.Error())
		return
	}

	apiReturn.Success(c)
}

// 创建并写入文件的函数
func createAndWriteFile(filename, content string) error {
	// 创建文件，如果文件已存在则会被覆盖
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 将内容写入文件
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	// fmt.Println("文件创建并写入成功:", filename)
	return nil
}

// 读取文件内容的函数
func readFile(filename string) (string, error) {
	// 读取文件内容
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
