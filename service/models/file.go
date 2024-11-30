package models

type File struct {
	BaseModel
	Src      string `json:"src"`
	UserId   uint   `json:"userId"`
	FileName string `json:"fileName" gorm:"varchar(255)"` // 文件名
	Type     int    `gorm:"tinyint(1)" json:"type"`       // 文件类型 1.壁纸 2.卡片图标
	Ext      string `gorm:"varchar(255)" json:"ext"`      // 扩展名
}

const (
	FileTypeUnknown = iota
	FileTypeWallpaper
	FileTypeCardIcon
)

// 添加一个文件记录
func (m *File) AddFile(userId uint, fileName, ext, src string, fileType int) (File, error) {
	file := File{
		UserId:   userId,
		FileName: fileName,
		Src:      src,
		Ext:      ext,
		Type:     fileType,
	}
	err := Db.Create(&file).Error

	return file, err
}
