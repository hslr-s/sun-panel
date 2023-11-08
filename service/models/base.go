package models

import (
	"fmt"
	"time"

	_ "gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// int类型代表是否的常量
const (
	INT_FALSE = iota
	INT_TURE
)

type BaseModel struct {
	gorm.Model
	// Db *gorm.DB `gorm:"_"`
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createTime"`
	UpdatedAt time.Time `json:"updateTime"`
	// DeletedAt DeletedAt `gorm:"index"`
}

type BaseModelNoId struct {
	CreatedAt time.Time      `json:"createTime"`
	UpdatedAt time.Time      `json:"updateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// 分页的结构体
type PageLimitStruct struct {
	PageSize  int `gorm:"-"` //
	LimitSize int `gorm:"-"` //
}

// 计算分页
func calcPage(page_size, limit_size int) (offset, limit int) {
	offset = limit_size * (page_size - 1)
	limit = limit_size
	return
}

var Db *gorm.DB

func GetDb() (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	dbModels := []interface{}{
		&User{},
	}
	dbDrive := "mysql"

	if dbDrive == "mysql" {

		db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(dbModels...)
		sqlDb, _ := db.DB()
		sqlDb.SetMaxIdleConns(10)             // SetMaxIdleConns 设置空闲连接池中连接的最大数量
		sqlDb.SetMaxOpenConns(100)            // SetMaxOpenConns 设置打开数据库连接的最大数量。
		sqlDb.SetConnMaxLifetime(time.Minute) // SetConnMaxLifetime 设置了连接可复用的最大时间。

	} else {
		fmt.Println("数据库驱动:", "SQLite")
		db, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				// TablePrefix:   "blog_",
				SingularTable: true,
			},
			// Logger: GetLogger(),
			// DisableForeignKeyConstraintWhenMigrating: true,
		})

		db.AutoMigrate(dbModels...)

	}

	Db = db
	return Db, err
}

// // 日志
// func GetLogger() logger.Interface {
// 	return logger.New(
// 		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
// 		logger.Config{
// 			SlowThreshold:             time.Second, // 慢 SQL 阈值
// 			LogLevel:                  logger.Warn, // 日志级别
// 			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
// 			Colorful:                  true,        // 彩色打印
// 		},
// 	)

// }
