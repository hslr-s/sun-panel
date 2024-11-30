package system

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"sun-panel/api/api_v1/common/apiData/commonApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type FileApi struct{}

func (a *FileApi) UploadImgOld(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	configUpload := global.Config.GetValueString("base", "source_path")

	// 文件类型
	fileTypeStr := c.PostForm("type")
	fileType, _ := strconv.Atoi(fileTypeStr)

	f, err := c.FormFile("imgfile")
	if err != nil {
		apiReturn.ErrorByCode(c, 1300)
		return
	} else {
		fileExt := strings.ToLower(path.Ext(f.Filename))
		agreeExts := []string{
			".png",
			".jpg",
			".gif",
			".jpeg",
			".webp",
			".svg",
			".ico",
			".avif",
			".heic",
		}

		if !cmn.InArray(agreeExts, fileExt) {
			apiReturn.ErrorByCode(c, 1301)
			return
		}
		fileName := cmn.Md5(fmt.Sprintf("%s%s", f.Filename, time.Now().String()))
		fildDir := fmt.Sprintf("%s/%d/%d/%d/", configUpload, time.Now().Year(), time.Now().Month(), time.Now().Day())
		isExist, _ := cmn.PathExists(fildDir)
		if !isExist {
			os.MkdirAll(fildDir, os.ModePerm)
		}
		filepath := fmt.Sprintf("%s%s%s", fildDir, fileName, fileExt)
		c.SaveUploadedFile(f, filepath)

		pureFilePath := filepath[len(configUpload):]
		// 向数据库添加记录
		mFile := models.File{}
		mFile.AddFile(userInfo.ID, f.Filename, fileExt, pureFilePath, fileType)
		apiReturn.SuccessData(c, gin.H{
			"imageUrl": global.UPLOAD_ROUTE + pureFilePath,
		})
	}
}

func (a *FileApi) UploadImg(c *gin.Context) {
	fileInfo, err := base.UploadFileAndSave(c, []string{
		".png",
		".jpg",
		".gif",
		".jpeg",
		".webp",
		".svg",
		".ico",
		".avif",
		".heic",
	}, 0)

	if err != nil {
		switch err {
		case base.ErrUploadExtensionNameNotAllowed:
			apiReturn.ErrorByCode(c, apiReturn.ErrCodeUnsupportedFileFormat)
		default:
			global.Logger.Error(err.Error())
			apiReturn.ErrorByCode(c, apiReturn.ErrCodeUploadFailed)
		}
		return
	}

	apiReturn.SuccessData(c, gin.H{
		"imageUrl": global.UPLOAD_ROUTE + fileInfo.Src,
	})
}

// 一直未使用有待测试
func (a *FileApi) UploadFiles(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	configUpload := global.Config.GetValueString("base", "source_path")

	// 文件类型
	fileTypeStr := c.PostForm("type")
	fileType, _ := strconv.Atoi(fileTypeStr)

	form, err := c.MultipartForm()
	if err != nil {
		apiReturn.ErrorByCode(c, 1300)
		return
	}
	files := form.File["files[]"]
	errFiles := []string{}
	succMap := map[string]string{}
	for _, f := range files {
		fileExt := strings.ToLower(path.Ext(f.Filename))
		fileName := cmn.Md5(fmt.Sprintf("%s%s", f.Filename, time.Now().String()))
		fildDir := fmt.Sprintf("%s/%d/%d/%d/", configUpload, time.Now().Year(), time.Now().Month(), time.Now().Day())
		isExist, _ := cmn.PathExists(fildDir)
		if !isExist {
			os.MkdirAll(fildDir, os.ModePerm)
		}
		filepath := fmt.Sprintf("%s%s%s", fildDir, fileName, fileExt)
		if c.SaveUploadedFile(f, filepath) != nil {
			errFiles = append(errFiles, f.Filename)
		} else {
			// 成功
			pureFilePath := filepath[len(configUpload):]
			// 向数据库添加记录
			mFile := models.File{}
			mFile.AddFile(userInfo.ID, f.Filename, fileExt, pureFilePath, fileType)
			succMap[f.Filename] = global.UPLOAD_ROUTE + pureFilePath
		}
	}

	apiReturn.SuccessData(c, gin.H{
		"succMap":  succMap,
		"errFiles": errFiles,
	})
}

func (a *FileApi) GetList(c *gin.Context) {
	list := []models.File{}
	userInfo, _ := base.GetCurrentUserInfo(c)
	var count int64
	if err := global.Db.Order("created_at desc").Find(&list, "user_id=?", userInfo.ID).Count(&count).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	data := []map[string]interface{}{}
	for _, v := range list {
		// 兼容旧版
		if len(v.Src) >= 0 && v.Src[0] == '.' {
			v.Src = v.Src[1:]
		} else {
			v.Src = global.UPLOAD_ROUTE + v.Src
		}
		data = append(data, map[string]interface{}{
			"src":        v.Src,
			"fileName":   v.FileName,
			"id":         v.ID,
			"createTime": v.CreatedAt,
			"updateTime": v.UpdatedAt,
			"path":       v.Src,
		})
	}
	apiReturn.SuccessListData(c, data, count)
}

func (a *FileApi) Deletes(c *gin.Context) {
	req := commonApiStructs.RequestDeleteIds[uint]{}
	userInfo, _ := base.GetCurrentUserInfo(c)
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	err := global.Db.Transaction(func(tx *gorm.DB) error {
		files := []models.File{}

		if err := tx.Order("created_at desc").Find(&files, "user_id=? AND id in ?", userInfo.ID, req.Ids).Error; err != nil {
			return err
		}

		for _, v := range files {
			// 兼容旧版
			var err error
			var path = ""
			if len(v.Src) >= 0 && v.Src[0] == '.' {
				path = v.Src
			} else {
				// 新版的储存方式
				sourcePath := global.Config.GetValueString("base", "source_path")
				path = strings.TrimRight(sourcePath, "/") + "/" + strings.TrimLeft(v.Src, "/")
			}

			// 判断文件是否存在
			exists := false
			if exists, err = cmn.PathExists(path); err == nil && exists {
				err = os.Remove(path)
			}

			if err != nil {
				global.Logger.Errorln("delete file ", path, ":", err)
				return err
			}
		}

		if err := tx.Order("created_at desc").Delete(&files, "user_id=? AND id in ?", userInfo.ID, req.Ids).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		apiReturn.Error(c, err.Error())
		return
	}
	apiReturn.Success(c)

}
