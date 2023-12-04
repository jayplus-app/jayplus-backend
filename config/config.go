package config

func LoadConfig() {
	LoadAppConfig()
	LoadDBConfig()
	LoadAuthConfig()
}
