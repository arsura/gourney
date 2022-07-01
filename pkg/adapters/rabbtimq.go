package adapter

import (
	config "github.com/arsura/gourney/configs"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type RabbitMQConnection struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      *amqp.Queue
}

func NewRabbitMQConnection(logger *zap.SugaredLogger, config *config.Config) *RabbitMQConnection {
	conn, err := amqp.Dial(config.RabbitMQ.URI)
	if err != nil {
		logger.With("error", err).Panic("failed to connect to rabbitmq")
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.With("error", err).Panic("failed to open a channel")
	}

	err = ch.ExchangeDeclare(
		config.RabbitMQ.Exchanges.LogsTopic,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.With("error", err).Panicf("failed to declare a %s exchange", config.RabbitMQ.Exchanges.LogsTopic)
	}

	return &RabbitMQConnection{
		Connection: conn,
		Channel:    ch,
	}
}
