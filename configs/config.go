package configs

import (
	"github.com/spf13/viper"
)

var cfg *config

type config struct {
	Log     Log
	Server  Server
	Swagger Swagger
}

type Log struct {
	Environment string
	Application string
}

type Server struct {
	Port string
}

type Swagger struct {
	Enabled bool
}

func init() {
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("LOG_ENVIRONMENT", "")
	viper.SetDefault("LOG_APPLICATION", "")
	viper.SetDefault("SWAGGER_ENABLED", false)

	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()
	viper.ReadInConfig()

	var la, le string
	if la = viper.GetString("LOG_APPLICATION"); la == "" {
		la = viper.GetString("APP_NAME")
	}
	if le = viper.GetString("LOG_ENVIRONMENT"); le == "" {
		le = viper.GetString("DOMAIN")
	}

	cfg = &config{
		Log: Log{
			Environment: le,
			Application: la,
		},
		Server: Server{
			Port: viper.GetString("PORT"),
		},
		Swagger: Swagger{
			Enabled: viper.GetBool("SWAGGER_ENABLED"),
		},
	}
}

func GetConfig() config {
	return *cfg
}
