package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type APIServer struct {
	Port     string
	IsEnable bool
}

type BlogCollections struct {
	Posts string
}

type BlogDatabase struct {
	Name        string
	Collections BlogCollections
}

type LogCollections struct {
	Logs string
}

type LogDatabase struct {
	Name        string
	Collections LogCollections
}

type MongoDB struct {
	URI string
	BlogDatabase
	LogDatabase
}

type Config struct {
	APIServer
	MongoDB
}

func NewConfig() *Config {
	env := os.Getenv("APP_ENV")

	switch env {
	case "development":
		viper.SetConfigName("local.config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("config")
	default:
		panic(fmt.Errorf("APP_ENV must not be undefined"))
	}

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("failed to read config file, %v", err))
	}

	return &Config{
		APIServer: APIServer{
			Port:     viper.GetString("service.api.port"),
			IsEnable: viper.GetBool("service.api.enable"),
		},
		MongoDB: MongoDB{
			URI: viper.GetString("database.mongodb.uri"),
			BlogDatabase: BlogDatabase{
				Name: viper.GetString("database.mongodb.databases.blog"),
				Collections: BlogCollections{
					Posts: viper.GetString("database.mongodb.collections.posts"),
				},
			},
			LogDatabase: LogDatabase{
				Name: viper.GetString("database.mongodb.databases.log"),
				Collections: LogCollections{
					Logs: viper.GetString("database.mongodb.collections.logs"),
				},
			},
		},
	}
}
