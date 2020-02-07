//TODO: To make prod ver - change ReadFile directory of config file to "config.json" and db host from mongo docker hostS to "localhost"
package main

import (
	"./config"
	"./server"
	"fmt"
)



func main() {
	cfg := config.GetConfig()
	app := &server.App{}
	app.Init(cfg)
	app.Run(":8080")
	fmt.Scanln()
}
