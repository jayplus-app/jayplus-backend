package config

import "os"

var (
	Port string
)

func LoadAppConfig() {
	Port = os.Getenv("SERVER_PORT")
}
