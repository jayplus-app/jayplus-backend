package config

import "os"

var (
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
)

func LoadDBConfig() {
	DBHost = os.Getenv("DB_HOST")
	DBPort = os.Getenv("DB_PORT")
	DBName = os.Getenv("DB_NAME")
	DBUser = os.Getenv("DB_USER")
	DBPassword = os.Getenv("DB_PASSWORD")
}
