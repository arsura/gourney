package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	PostgresURI string
	APIServer
}

type APIServer struct {
	Port     string
	IsEnable bool
}

func NewConfig() *Config {
	env := os.Getenv("APP_ENV")

	switch env {
	case "development":
		viper.SetConfigName("local.config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("config")
	default:
		panic(fmt.Errorf("env is undefined"))
	}

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %v", err))
	}

	return &Config{
		PostgresURI: viper.GetString("database.postgres.uri"),
		APIServer: APIServer{
			Port:     viper.GetString("service.api.port"),
			IsEnable: viper.GetBool("service.api.is_enable"),
		},
	}
}
