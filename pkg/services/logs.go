package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/arsura/gourney/config"
	adapter "github.com/arsura/gourney/pkg/adapters"
	"github.com/arsura/gourney/pkg/constant"
	model "github.com/arsura/gourney/pkg/models/mongodb"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type LogEventBody struct {
	Id       primitive.ObjectID   `json:"id"`
	Event    constant.EventAction `json:"event"`
	ActionAt time.Time            `json:"action_at"`
}

type LogServiceProvider interface {
	PublishPostLogEvent(ctx context.Context, event constant.EventAction, post *model.Post) error
}

type logService struct {
	rabbitmqConnection *adapter.RabbitMQConnection
	logger             *zap.SugaredLogger
	config             *config.Config
}

func NewLogService(rabbitmqConnection *adapter.RabbitMQConnection, logger *zap.SugaredLogger, config *config.Config) *logService {
	return &logService{rabbitmqConnection, logger, config}
}

func (s *logService) PublishPostLogEvent(ctx context.Context, event constant.EventAction, post *model.Post) error {
	body := &LogEventBody{
		Id:       post.Id,
		Event:    event,
		ActionAt: time.Now(),
	}
	bodyAsJson, err := json.Marshal(body)
	if err != nil {
		s.logger.With("error", err).Error("failed to marshal body to json")
		return err
	}

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        bodyAsJson,
		MessageId:   ctx.Value(constant.REQUEST_ID_KEY).(string),
	}

	routingKey := fmt.Sprintf("%s.%s.%s", constant.ROUTING_KEY_PREFIX_LOGS, constant.TARGET_POST, event)
	err = s.rabbitmqConnection.Channel.Publish(
		s.config.RabbitMQ.Exchanges.LogsTopic,
		routingKey,
		false,
		false,
		msg,
	)
	if err != nil {
		s.logger.With("error", err, "routing_key", routingKey, "message", msg).Error("failed to publish post log event")
		return err
	}

	s.logger.With("message", msg, "routing_key", routingKey).Info("publish post changes event success")
	return nil
}
