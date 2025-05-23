package config

import (
	"tugas-akhir/drivers/database"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if os.Getenv("NAME") != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func InitConfigDB() database.Config {
	return database.Config{
		DB_HOST:      os.Getenv("DB_HOST"),
		DB_USERNAME:  os.Getenv("DB_USERNAME"),
		DB_PASSWORD:  os.Getenv("DB_PASSWORD"),
		DB_NAME:      os.Getenv("DB_NAME"),
		DB_PORT:      os.Getenv("DB_PORT"),
		DB_SSL:       os.Getenv("DB_SSL"),
		DB_TZ:        os.Getenv("DB_TZ"),
		DB_LOG_LEVEL: os.Getenv("DB_LOG_LEVEL"),
	}
}

