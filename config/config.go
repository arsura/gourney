package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type LogsWorkerService struct {
	IsEnable bool
}

type APIService struct {
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
	PostLogs string
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

type Exchanges struct {
	LogsTopic string
}

type RabbitMQ struct {
	URI string
	Exchanges
}

type Config struct {
	APIService
	LogsWorkerService
	MongoDB
	RabbitMQ
}

func NewConfig(logger *zap.SugaredLogger) *Config {
	viper.AutomaticEnv()

	err := viper.BindEnv("service.api.enable", "ENABLE_API")
	if err != nil {
		logger.With("error", err).Panic("failed to bind ENABLE_API env")
	}

	err = viper.BindEnv("service.worker.enable", "ENABLE_WORKER")
	if err != nil {
		logger.With("error", err).Panic("failed to bind ENABLE_WORKER env")
	}

	viper.SetConfigName("local.config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")

	err = viper.ReadInConfig()
	if err != nil {
		logger.With("error", err).Panic("failed to read config file")
	}

	return &Config{
		APIService: APIService{
			Port:     viper.GetString("service.api.port"),
			IsEnable: viper.GetBool("service.api.enable"),
		},
		LogsWorkerService: LogsWorkerService{
			IsEnable: viper.GetBool("service.worker.enable"),
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
					PostLogs: viper.GetString("database.mongodb.collections.post_logs"),
				},
			},
		},
		RabbitMQ: RabbitMQ{
			URI: viper.GetString("broker.rabbitmq.uri"),
			Exchanges: Exchanges{
				LogsTopic: viper.GetString("broker.rabbitmq.exchanges.logs_topic"),
			},
		},
	}
}
