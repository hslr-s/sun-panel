package biz

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/lib/siteFavicon"
	"sun-panel/models"
	"time"

	"gorm.io/gorm"
)

type FileType struct {
}

// 下载图片并将记录保存到数据库中
func (p *FileType) DownloadFileAndSaveDb(db *gorm.DB, fileUrl string, userId uint) (fileInfo models.File, err error) {

	parsedURL, err := url.Parse(fileUrl)
	if err != nil {
		global.Logger.Errorln("error: URL parsing failed("+fileUrl+"):", err.Error())
		return
	}

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
		if imgInfo, err = siteFavicon.DownloadImage(fileUrl, savePath, 1024*1024); err != nil {
			return
		}
	}

	pureFilePath := imgInfo.Name()[len(configUpload):]
	ext := path.Ext(fileUrl)

	// 保存到数据库
	fileInfo = models.File{
		UserId:   userId,
		FileName: parsedURL.Host,
		Src:      pureFilePath,
		Type:     models.FileTypeCardIcon,
		Ext:      ext,
	}
	err = db.Create(&fileInfo).Error
	return
}
