package main

import (
	"./config"
	"./server"
	"fmt"
)

func main() {
	config := config.GetConfig()
	app := &server.App{}
	app.Init(config)
	app.Run(":8080")
	fmt.Scanln()
}