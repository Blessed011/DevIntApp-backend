package dsn

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// FromEnv собирает DSN строку из переменных окружения
func FromEnv() string {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Print("No .env file found")
	}
	host := os.Getenv("DB_HOST")
	if host == "" {
		return ""
	}

	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")
	fmt.Println(pass)
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
}
