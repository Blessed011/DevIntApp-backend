package main

import (
	"lab1/internal/pkg/app"
	"log"
)

// TODO: change
// @title BITOP
// @version 1.0
// @description Bmstu Open IT Platform

// @license.name AS IS (NO WARRANTY)

// @host 127.0.0.1
// @schemes https http
// @BasePath /

func main() {
	app, err := app.New()
	if err != nil {
		log.Println("app can not be created", err)
		return
	}
	app.Run()
}
