package main

import (
	"ITLabReports/config"
	"ITLabReports/server"
	"fmt"
)

func main() {
	cfg := config.GetConfig()
	app := &server.App{}
	app.Init(cfg)
	app.Run(":8080")
	fmt.Scanln()
}
