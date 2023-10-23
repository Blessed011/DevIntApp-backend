package main

import (
	"lab1/internal/app/ds"
	"lab1/internal/app/dsn"

	// "log"

	// "github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// if err := godotenv.Load("../../.env"); err != nil {
	// 	log.Print("No .env file found")
	// }
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&ds.Module{},
		&ds.Mission{},
		&ds.User{},
		&ds.Flight{},
	)
	if err != nil {
		panic(err)

	}
}
