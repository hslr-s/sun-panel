package initialize

import (
	"flag"
	"fmt"
	"os"
	"sun-panel/global"
	"sun-panel/initialize/cUserToken"
	"sun-panel/initialize/config"
	"sun-panel/initialize/database"
	"sun-panel/initialize/lang"
	"sun-panel/initialize/other"
	"sun-panel/initialize/redis"
	"sun-panel/initialize/runlog"
	"sun-panel/initialize/systemSettingCache"
	"sun-panel/initialize/userToken"
	"sun-panel/lib/cmn"
	"sun-panel/models"
	"sun-panel/structs"

	"log"

	"github.com/gin-gonic/gin"
)

var DB_DRIVER = database.SQLITE

// var RUNCODE = "debug"
// var ISDOCER = "" // 是否为docker模式

func InitApp() error {
	Logo()
	gin.SetMode(global.RUNCODE) // GIN 运行模式

	// 日志
	if logger, err := runlog.InitRunlog(global.RUNCODE, "running.log"); err != nil {
		log.Panicln("日志初始化错误", err)
		panic(err)
	} else {
		global.Logger = logger
	}

	// 命令行运行
	CommandRun()

	// 配置初始化
	{
		if config, err := config.ConfigInit(); err != nil {
			global.Logger.Errorln("配置初始化错误", err)
			return err
		} else {
			global.Config = config
		}
	}

	// 多语言初始化
	lang.LangInit("zh-cn") // en-us

	DatabaseConnect()

	// Redis 连接
	{
		// 判断是否有使用redis的驱动，没有将不连接
		cacheDrive := global.Config.GetValueString("base", "cache_drive")
		queueDrive := global.Config.GetValueString("base", "queue_drive")
		if cacheDrive == "redis" || queueDrive == "redis" {
			redisConfig := structs.IniConfigRedis{}
			global.Config.GetSection("redis", &redisConfig)
			rdb, err := redis.InitRedis(redis.Options{
				Addr:     redisConfig.Address,
				Password: redisConfig.Password,
				DB:       redisConfig.Db,
			})

			if err != nil {
				log.Panicln("Redis初始化错误", err)
				panic(err)
				// return err
			}
			global.RedisDb = rdb
		}
	}

	// 初始化用户token
	global.UserToken = userToken.InitUserToken()
	global.CUserToken = cUserToken.InitCUserToken()

	// 其他的初始化
	global.VerifyCodeCachePool = other.InitVerifyCodeCachePool()
	global.SystemSetting = systemSettingCache.InItSystemSettingCache()

	return nil
}

func DatabaseConnect() {
	// 数据库连接 - 开始
	var dbClientInfo database.DbClient
	databaseDrive := global.Config.GetValueStringOrDefault("base", "database_drive")
	if databaseDrive == database.MYSQL {
		dbClientInfo = &database.MySQLConfig{
			Username:    global.Config.GetValueStringOrDefault("mysql", "username"),
			Password:    global.Config.GetValueStringOrDefault("mysql", "password"),
			Host:        global.Config.GetValueStringOrDefault("mysql", "host"),
			Port:        global.Config.GetValueStringOrDefault("mysql", "port"),
			Database:    global.Config.GetValueStringOrDefault("mysql", "db_name"),
			WaitTimeout: global.Config.GetValueInt("mysql", "wait_timeout"),
		}
	} else {
		dbClientInfo = &database.SQLiteConfig{
			Filename: global.Config.GetValueStringOrDefault("sqlite", "file_path"),
		}
	}

	if db, err := database.DbInit(dbClientInfo); err != nil {
		log.Panicln("数据库初始化错误", err)
		panic(err)
	} else {
		global.Db = db
		models.Db = global.Db
	}

	database.CreateDatabase(databaseDrive, global.Db)

	database.NotFoundAndCreateUser(global.Db)
}

// 命令行运行
func CommandRun() {
	var (
		cfg bool
		pwd bool
	)

	flag.BoolVar(&cfg, "config", false, "生成配置文件")
	flag.BoolVar(&pwd, "password-reset", false, "重置第一个用户的密码")

	flag.Parse()

	if cfg {
		// 生成配置文件
		fmt.Println("正在生成配置文件")
		cmn.AssetsTakeFileToPath("conf.example.ini", "conf/conf.example.ini")
		cmn.AssetsTakeFileToPath("conf.example.ini", "conf/conf.ini")
		fmt.Println("配置文件已经创建 conf/conf.ini ", "请按照自己的需求修改")
		os.Exit(0) // 务必退出
	} else if pwd {
		// 重置密码

		// 配置初始化
		config, _ := config.ConfigInit()
		global.Config = config

		DatabaseConnect()
		userInfo := models.User{}
		if err := global.Db.Where("role=?", 1).Order("id").First(&userInfo).Error; err != nil {
			fmt.Println("ERROR", err.Error())
			os.Exit(0) // 务必退出
		}

		newPassword := "12345678"

		updateInfo := models.User{
			Password: cmn.PasswordEncryption(newPassword),
			Token:    "",
		}
		// 重置第一个管理员的密码
		if err := global.Db.Select("Password", "Token").Where("id=?", userInfo.ID).Updates(&updateInfo).Error; err != nil {
			fmt.Println("ERROR", err.Error())
			os.Exit(0) // 务必退出
		}

		fmt.Println("密码已经重置成功，以下是账号信息")
		fmt.Println("用户名 ", userInfo.Username)
		fmt.Println("密码 ", newPassword)
		os.Exit(0) // 务必退出
	} else {
		return
	}
	os.Exit(0) // 务必退出
}

func Logo() {
	fmt.Println("     ____            ___                __")
	fmt.Println("    / __/_ _____    / _ \\___ ____  ___ / /")
	fmt.Println("   _\\ \\/ // / _ \\  / ___/ _ `/ _ \\/ -_) / ")
	fmt.Println("  /___/\\_,_/_//_/ /_/   \\_,_/_//_/\\__/_/  ")
	fmt.Println("")

	versionInfo := cmn.GetSysVersionInfo()
	fmt.Println("Version:", versionInfo.Version)
	fmt.Println("Welcome to the Sun-Panel.")
	fmt.Println("Project address:", "https://github.com/hslr-s/sun-panel")

}
