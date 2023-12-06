package main

import (
	"lab1/internal/pkg/app"
	"log"
)

// TODO: change
// @title Lunar Gateway
// @version 1.0

// @host 127.0.0.1:8081
// @schemes http
// @BasePath /

func main() {
	app, err := app.New()
	if err != nil {
		log.Println("app can not be created", err)
		return
	}
	app.Run()
}
