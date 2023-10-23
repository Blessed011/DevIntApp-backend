package main

import (
	"lab1/internal/pkg/app"
	"log"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Println("app can not be created", err)
		return
	}
	app.Run()
}
