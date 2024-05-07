package main

import (
	"log"
	"sun-panel/global"
	"sun-panel/initialize"
	"sun-panel/router"
)

func main() {
	err := initialize.InitApp()
	if err != nil {
		log.Println("初始化错误:", err.Error())
		panic(err)
	}
	httpPort := global.Config.GetValueStringOrDefault("base", "http_port")
	httpAddress := global.Config.GetValueStringOrDefault("base", "http_address")
	if err := router.InitRouters(httpAddress + ":" + httpPort); err != nil {
		panic(err)
	}
}
