package panel

import (
	"os"
	"path"
	"strings"
	"sun-panel/api/api_v1/common/apiData/panelApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/lib/cmn"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type GlobalSettingApi struct {
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

	biz.GlobalSetting.GetSiteStting(&resp)
	apiReturn.SuccessData(c, resp)
}

// 获取全局站点设置
func (a *GlobalSettingApi) UploadFaviconImage(c *gin.Context) {
	customPath := global.Config.GetValueString("base", "custom_style_path")
	f, err := c.FormFile("imgfile")
	if err != nil {
		apiReturn.ErrorByCode(c, 1300)
		return
	} else {
		fileExt := strings.ToLower(path.Ext(f.Filename))
		agreeExts := []string{
			// ".png",
			// ".jpg",
			// ".gif",
			// ".jpeg",
			// ".webp",
			// ".svg",
			".ico",
		}

		if !cmn.InArray(agreeExts, fileExt) {
			apiReturn.ErrorByCode(c, 1301)
			return
		}

		isExist, _ := cmn.PathExists(customPath)
		if !isExist {
			os.MkdirAll(customPath, os.ModePerm)
		}
		filepath := customPath + "/favicon.ico"
		c.SaveUploadedFile(f, filepath)

		apiReturn.SuccessData(c, gin.H{
			"imageUrl": filepath,
		})
	}
}

func (a *GlobalSettingApi) GetCustomStylePath(c *gin.Context) {
	customPath := global.Config.GetValueString("base", "custom_style_path")
	apiReturn.SuccessData(c, gin.H{
		"customPath": customPath,
	})
}

func (a *GlobalSettingApi) GetCustomJsAndCssCode(c *gin.Context) {
	customPath := global.Config.GetValueString("base", "custom_style_path")

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
	customPath := global.Config.GetValueString("base", "custom_style_path")

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
		global.Logger.Error("createAndWriteFile:", cssFilePath, "-", err.Error())
		return
	}

	jsFilePath := customPath + "/index.js"
	// cssIsExist, _ := cmn.PathExists(customPath)
	if err := createAndWriteFile(jsFilePath, req.JsContent); err != nil {
		apiReturn.ErrorByCode(c, apiReturn.ErrCodeCustomJSAndCSSFilePreservationFailure)
		global.Logger.Error("createAndWriteFile:", cssFilePath, "-", err.Error())
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
