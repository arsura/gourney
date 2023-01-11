package logsworker

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/arsura/gourney/pkg/constant"
	model "github.com/arsura/gourney/pkg/models"
	service "github.com/arsura/gourney/pkg/services"
	"github.com/streadway/amqp"
)

func (app *Application) classifyTarget(msgs []amqp.Delivery) map[constant.TargetEvent][]amqp.Delivery {
	classifier := make(map[constant.TargetEvent][]amqp.Delivery)
	for _, msg := range msgs {
		target := constant.TargetEvent(strings.Split(msg.RoutingKey, ".")[1])
		switch target {
		case constant.TARGET_POST:
			classifier[target] = append(classifier[target], msg)
		default:
			app.Logger.With("event", "classify_target", "tracking_id", msg.MessageId, "message", msg).Warn("cannot classify message target")
		}
	}
	return classifier
}

func (app *Application) handler(msgs []amqp.Delivery) {
	classifier := app.classifyTarget(msgs)
	for target, msgs := range classifier {
		switch target {
		case constant.TARGET_POST:
			var postLogs []model.PostLog
			for _, msg := range msgs {
				logBody := &service.LogEventBody{}
				err := json.Unmarshal(msg.Body, logBody)
				if err != nil {
					app.Logger.With("error", err, "event", "handle_logs_message", "tracking_id", msg.MessageId).Error("failed to unmarshal message body")
					continue
				}
				postLog := model.PostLog{
					PostId:   logBody.Id,
					Event:    logBody.Event,
					ActionAt: logBody.ActionAt,
				}
				postLogs = append(postLogs, postLog)
			}
			app.Usecases.PostLog.CreatePostLogs(context.Background(), postLogs)
		default:
		}
	}
}
