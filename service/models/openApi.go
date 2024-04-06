package models

type OpenApiToken struct {
	BaseModel
	Token  string `gorm:"unique;type:varchar(255)"`
	UserId uint
	User   User
}
