package adapter

import (
	config "github.com/arsura/gourney/configs"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type Queues struct {
	Hello amqp.Queue
}

type RabbitMQConnection struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queues     *Queues
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

	helloQueue, err := ch.QueueDeclare(
		config.RabbitMQ.Queues.Hello, // name
		false,                        // durable
		false,                        // delete when unused
		false,                        // exclusive
		false,                        // no-wait
		nil,                          // arguments
	)
	if err != nil {
		logger.With("error", err).Panicf("failed to declare a %s queue", config.RabbitMQ.Queues.Hello)
	}

	return &RabbitMQConnection{
		Connection: conn,
		Channel:    ch,
		Queues: &Queues{
			Hello: helloQueue,
		},
	}
}
