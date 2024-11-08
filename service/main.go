package main

import (
	"log"
	"sun-panel/initialize"
	"sun-panel/router"
)

func main() {
	err := initialize.InitApp()
	if err != nil {
		log.Fatal("Initialization error:", err)
	}

	router.InitRouters()
}
