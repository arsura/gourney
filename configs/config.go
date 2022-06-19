package config

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
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

type Queues struct {
	Hello string
}

type RabbitMQ struct {
	URI string
	Queues
}

type Config struct {
	APIServer
	MongoDB
	RabbitMQ
}

func NewConfig(logger *zap.SugaredLogger) *Config {
	env := os.Getenv("APP_ENV")

	switch env {
	case "development":
		viper.SetConfigName("local.config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("configs")
	default:
		logger.Panic("APP_ENV must not be undefined")
	}

	err := viper.ReadInConfig()
	if err != nil {
		logger.With("error", err).Panic("failed to read config file")
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
		RabbitMQ: RabbitMQ{
			URI: viper.GetString("broker.rabbitmq.uri"),
			Queues: Queues{
				Hello: viper.GetString("broker.rabbitmq.queues.hello"),
			},
		},
	}
}
