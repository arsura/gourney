package logsworker

import (
	"time"

	"github.com/arsura/gourney/config"
	adapter "github.com/arsura/gourney/pkg/adapters"
	"github.com/arsura/gourney/pkg/constant"
	usecase "github.com/arsura/gourney/pkg/usecases"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type Application struct {
	RabbitMQConnection *adapter.RabbitMQConnection
	Usecases           *usecase.Usecase
	Logger             *zap.SugaredLogger
	Config             *config.Config
}

func NewWorkerApplication(rabbitMQConnection *adapter.RabbitMQConnection, usecases *usecase.Usecase, logger *zap.SugaredLogger, config *config.Config) *Application {
	return &Application{rabbitMQConnection, usecases, logger, config}
}

func (app *Application) Start() {
	q, err := app.RabbitMQConnection.Channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		app.Logger.With("error", err).Panic("failed to declare a queue")
	}

	err = app.RabbitMQConnection.Channel.QueueBind(
		q.Name,
		constant.ROUTING_KEY_LOGS,
		app.Config.RabbitMQ.Exchanges.LogsTopic,
		false,
		nil)
	if err != nil {
		app.Logger.With("error", err).Panic("failed to bind a queue")
	}

	msgs, err := app.RabbitMQConnection.Channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		app.Logger.With("error", err).Panic("failed to register a consumer")
	}

	var (
		forever chan struct{}
		tmpMsgs []amqp.Delivery

		ticker = time.NewTicker(constant.TICK_TIME_TO_WRITE_LOGS)
	)

	go func() {
		for {
			select {
			case msg := <-msgs:
				app.Logger.With("message", msg).Info("receive message")
				tmpMsgs = append(tmpMsgs, msg)
				if len(tmpMsgs) >= constant.MAX_TEMP_MESSAGE_SIZE {
					app.Logger.Info("the temporary buffer had reached its limit, write logs")
					app.handler(tmpMsgs)
					tmpMsgs = nil
				}
			case <-ticker.C:
				app.Logger.Infof("receive tick signal to write %d message of logs", len(tmpMsgs))
				app.handler(tmpMsgs)
				tmpMsgs = nil
			}
		}
	}()

	app.Logger.Info("start worker service success, waiting for messages. to exit press ctrl+c")
	<-forever
}
