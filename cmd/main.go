package main

import (
	"nsparser/internal/app"
	"nsparser/internal/config"
	"os"
)

func main() {
	path, ok := os.LookupEnv("NS_DL_CONF_PATH")
	if !ok {
		path = "config.json"
	}
	c := config.NewConfig(path)
	defer c.Save()
	app := app.Init(c)
	app.Run()
}
