package panel

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sun-panel/api/api_v1/common/apiData/commonApiStructs"
	"sun-panel/api/api_v1/common/apiData/panelApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/lib/siteFavicon"
	"sun-panel/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type ItemIcon struct {
}

func (a *ItemIcon) Edit(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	req := panelApiStructs.ItemInfo{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	req.UserId = userInfo.ID

	saveData := models.ItemIcon{}
	saveData.ID = req.ID
	saveData.CardType = req.CardType
	saveData.IconJson = req.IconJson
	saveData.Icon = req.Icon
	saveData.Title = req.Title
	saveData.OnlyName = req.OnlyName
	saveData.Url = req.Url
	saveData.LanUrl = req.LanUrl
	saveData.Description = req.Description
	saveData.OpenMethod = req.OpenMethod
	saveData.Sort = req.Sort
	saveData.BackgroundColor = req.BackgroundColor
	saveData.ItemIconGroupId = req.ItemIconGroupId
	saveData.ExpandParam = base.AnyToJsonStr(req.ExpandParam)
	saveData.UserId = userInfo.ID

	if v, err := biz.Item.Edit(userInfo.ID, saveData); err != nil {
		if errors.Is(err, biz.ErrItemOnlyNameExisted) {
			apiReturn.ErrorByCodeAndMsg(c, -2, err.Error())
			return
		}
		apiReturn.ErrorDatabase(c, err.Error())
		return
	} else {
		saveData = v
	}

	apiReturn.SuccessData(c, saveData)
}

// 添加多个图标
func (a *ItemIcon) AddMultiple(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	// type Request
	req := []panelApiStructs.ItemInfo{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	datas := []models.ItemIcon{}

	for _, v := range req {
		if v.ItemIconGroupId == 0 {
			apiReturn.ErrorParamFomat(c, "Group is mandatory")
			return
		}
		v.UserId = userInfo.ID
		// json转字符串
		if j, err := json.Marshal(v.Icon); err == nil {
			v.IconJson = string(j)
		}

		saveData := models.ItemIcon{}

		saveData.CardType = v.CardType
		saveData.IconJson = v.IconJson
		saveData.Icon = v.Icon
		saveData.Title = v.Title
		saveData.OnlyName = v.OnlyName
		saveData.Url = v.Url
		saveData.LanUrl = v.LanUrl
		saveData.Description = v.Description
		saveData.OpenMethod = v.OpenMethod
		saveData.Sort = v.Sort
		saveData.BackgroundColor = v.BackgroundColor
		saveData.ItemIconGroupId = v.ItemIconGroupId
		saveData.ExpandParam = base.AnyToJsonStr(v.ExpandParam)
		saveData.UserId = userInfo.ID
		datas = append(datas, saveData)
	}

	global.Db.Create(&datas)

	apiReturn.SuccessData(c, req)
}

// // 获取详情
// func (a *ItemIcon) GetInfo(c *gin.Context) {
// 	req := systemApiStructs.AiDrawGetInfoReq{}

// 	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
// 		apiReturn.ErrorParamFomat(c, err.Error())
// 		return
// 	}

// 	userInfo, _ := base.GetCurrentUserInfo(c)

// 	aiDraw := models.AiDraw{}
// 	aiDraw.ID = req.ID
// 	if err := aiDraw.GetInfo(global.Db); err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			apiReturn.Error(c, "不存在记录")
// 			return
// 		}
// 		apiReturn.ErrorDatabase(c, err.Error())
// 		return
// 	}

// 	if userInfo.ID != aiDraw.UserID {
// 		apiReturn.ErrorNoAccess(c)
// 		return
// 	}

// 	apiReturn.SuccessData(c, aiDraw)
// }

func (a *ItemIcon) GetListByGroupId(c *gin.Context) {
	req := models.ItemIcon{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)
	itemIcons := []models.ItemIcon{}

	resp := []panelApiStructs.ItemInfo{}

	if err := global.Db.Order("sort ,created_at").Find(&itemIcons, "item_icon_group_id = ? AND user_id=?", req.ItemIconGroupId, userInfo.ID).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	for _, v := range itemIcons {
		json.Unmarshal([]byte(v.IconJson), &v.Icon)
		temp := panelApiStructs.ItemInfo{}

		temp.ID = v.ID
		temp.CardType = v.CardType
		temp.Icon = v.Icon
		temp.Title = v.Title
		temp.OnlyName = v.OnlyName
		temp.Url = v.Url
		temp.LanUrl = v.LanUrl
		temp.Description = v.Description
		temp.OpenMethod = v.OpenMethod
		temp.Sort = v.Sort
		temp.BackgroundColor = v.BackgroundColor
		temp.ItemIconGroupId = v.ItemIconGroupId
		temp.UserId = userInfo.ID
		temp.ExpandParam = base.JsonStrToMap(v.ExpandParam)

		resp = append(resp, temp)
	}

	apiReturn.SuccessListData(c, resp, 0)
}

func (a *ItemIcon) Deletes(c *gin.Context) {
	req := commonApiStructs.RequestDeleteIds[uint]{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)
	if err := global.Db.Delete(&models.ItemIcon{}, "id in ? AND user_id=?", req.Ids, userInfo.ID).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

// 保存排序
func (a *ItemIcon) SaveSort(c *gin.Context) {
	req := panelApiStructs.ItemIconSaveSortRequest{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)

	transactionErr := global.Db.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		for _, v := range req.SortItems {
			if err := tx.Model(&models.ItemIcon{}).Where("user_id=? AND id=? AND item_icon_group_id=?", userInfo.ID, v.Id, req.ItemIconGroupId).Update("sort", v.Sort).Error; err != nil {
				// 返回任何错误都会回滚事务
				return err
			}
		}

		// 返回 nil 提交事务
		return nil
	})

	if transactionErr != nil {
		apiReturn.ErrorDatabase(c, transactionErr.Error())
		return
	}

	apiReturn.Success(c)
}

// 支持获取并直接下载对方网站图标到服务器[将被弃用] 2024年11月30日修改
func (a *ItemIcon) GetSiteFavicon(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	req := panelApiStructs.ItemIconGetSiteFaviconReq{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	resp := panelApiStructs.ItemIconGetSiteFaviconResp{}
	fullUrl := ""
	if iconUrl, err := siteFavicon.GetOneFaviconURL(req.Url); err != nil {
		// 尝试使用最原始的方式获取图标 https://域名/favicon.ico
		if urlInfo, err := url.Parse(req.Url); err != nil {
			errMsg := "acquisition failed: get ico error:" + err.Error()
			global.Logger.Error(errMsg)
			apiReturn.Error(c, errMsg)
			return
		} else {
			fullUrl = urlInfo.Scheme + "://" + urlInfo.Host + "/favicon.ico"
		}

	} else {
		fullUrl = iconUrl
	}

	parsedURL, err := url.Parse(req.Url)
	if err != nil {
		apiReturn.Error(c, "acquisition failed:"+err.Error())
		return
	}

	protocol := parsedURL.Scheme
	global.Logger.Debug("protocol:", protocol)
	global.Logger.Debug("fullUrl:", fullUrl)

	// 如果URL以双斜杠（//）开头，则使用当前页面协议
	if strings.HasPrefix(fullUrl, "//") {
		fullUrl = protocol + "://" + fullUrl[2:]
	} else if !strings.HasPrefix(fullUrl, "http://") && !strings.HasPrefix(fullUrl, "https://") {
		// 如果URL既不以http://开头也不以https://开头，则默认为http协议
		fullUrl = "http://" + fullUrl
	}
	global.Logger.Debug("fullUrl:", fullUrl)
	// 去除图标的get参数
	{
		parsedIcoURL, err := url.Parse(fullUrl)
		if err != nil {
			apiReturn.Error(c, "acquisition failed: parsed ico URL :"+err.Error())
			return
		}
		fullUrl = parsedIcoURL.Scheme + "://" + parsedIcoURL.Host + parsedIcoURL.Path
	}
	global.Logger.Debug("fullUrl:", fullUrl)

	// 生成保存目录
	configUpload := global.Config.GetValueString("base", "source_path")
	savePath := fmt.Sprintf("%s/%d/%d/%d/", configUpload, time.Now().Year(), time.Now().Month(), time.Now().Day())
	global.Logger.Debug("保存路径", savePath)
	isExist, _ := cmn.PathExists(savePath)
	if !isExist {
		os.MkdirAll(savePath, os.ModePerm)
	}

	// 下载
	var imgInfo *os.File
	{
		var err error
		if imgInfo, err = siteFavicon.DownloadImage(fullUrl, savePath, 1024*1024); err != nil {
			apiReturn.Error(c, "acquisition failed: download"+err.Error())
			return
		}
	}

	pureFilePath := imgInfo.Name()[len(configUpload):]
	// 保存到数据库
	ext := path.Ext(fullUrl)
	mFile := models.File{}
	if _, err := mFile.AddFile(userInfo.ID, parsedURL.Host, ext, pureFilePath, models.FileTypeCardIcon); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	resp.IconUrl = global.UPLOAD_ROUTE + pureFilePath
	apiReturn.SuccessData(c, resp)
}

func (a *ItemIcon) GetSiteFaviconList(c *gin.Context) {
	// userInfo, _ := base.GetCurrentUserInfo(c)
	req := panelApiStructs.ItemIconGetSiteFaviconReq{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	fullUrls := []string{}
	if iconUrls, err := siteFavicon.GetFaviconURLs(req.Url); err != nil {
		// 尝试使用最原始的方式获取图标 https://域名/favicon.ico
		if urlInfo, err := url.Parse(req.Url); err != nil {
			errMsg := "acquisition failed: get ico error:" + err.Error()
			global.Logger.Error(errMsg)
			apiReturn.Error(c, errMsg)
			return
		} else {
			fullUrls = append(fullUrls, urlInfo.Scheme+"://"+urlInfo.Host+"/favicon.ico")
		}

	} else {
		fullUrls = iconUrls
	}

	parsedURL, err := url.Parse(req.Url)
	if err != nil {
		apiReturn.Error(c, "acquisition failed:"+err.Error())
		return
	}

	protocol := parsedURL.Scheme
	global.Logger.Debug("protocol:", protocol)
	global.Logger.Debug("fullUrls:", fullUrls)

	// 调整地址
	for key, fullUrl := range fullUrls {
		// 如果URL以双斜杠（//）开头，则使用当前页面协议
		if strings.HasPrefix(fullUrl, "//") {
			fullUrl = protocol + "://" + fullUrl[2:]
		} else if !strings.HasPrefix(fullUrl, "http://") && !strings.HasPrefix(fullUrl, "https://") {
			// 如果URL既不以http://开头也不以https://开头，则默认为http协议
			fullUrls[key] = "http://" + fullUrl
		}

		// 去除图标的get参数
		{
			parsedIcoURL, err := url.Parse(fullUrl)
			if err != nil {
				apiReturn.Error(c, "acquisition failed: parsed ico URL :"+err.Error())
				return
			}
			fullUrls[key] = parsedIcoURL.Scheme + "://" + parsedIcoURL.Host + parsedIcoURL.Path
		}

		global.Logger.Debug("fullUrl:", fullUrl)
	}

	apiReturn.SuccessListData(c, fullUrls, int64(len(fullUrls)))
}

// 图片代理
func (a *ItemIcon) ImageProxy(c *gin.Context) {
	imageUrl := c.Request.URL.Query().Get("imageUrl")
	if imageUrl == "" {
		http.Error(c.Writer, "imageUrl parameter is required", http.StatusBadRequest)
		return
	}

	resp, err := http.Get(imageUrl)
	if err != nil {
		http.Error(c.Writer, "Failed to fetch image", http.StatusInternalServerError)
		log.Printf("Error fetching image: %v", err)
		return
	}
	defer resp.Body.Close()

	// Copy the response body to our response writer
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		http.Error(c.Writer, "Failed to copy image", http.StatusInternalServerError)
		log.Printf("Error copying image: %v", err)
		return
	}
}

// 下载图标图像
func (a *ItemIcon) DownloadFavicon(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	req := panelApiStructs.ItemIconGetSiteFaviconReq{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	info, err := base.DownloadFileAndSave(req.Url, userInfo.ID)
	if err != nil {
		apiReturn.Error(c, err.Error())
		return
	}

	resp := panelApiStructs.ItemIconGetSiteFaviconResp{}
	resp.IconUrl = global.UPLOAD_ROUTE + info.Src
	// global.Logger.Debug("download:", resp.IconUrl)
	apiReturn.SuccessData(c, resp)
}
