//TODO: To make prod ver - change ReadFile directory of config file to "congif.json" and db host from mongo docker hostS to "localhost"
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
