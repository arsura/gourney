package logger

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

func NewLogger() *zap.SugaredLogger {
	rawJSON := []byte(`{
		"level": "debug",
		"encoding": "json",
		"outputPaths": ["stdout", "/app/logs"],
		"errorOutputPaths": ["stderr"],
		"encoderConfig": {
		  "messageKey": "message",
		  "levelKey": "level",
		  "levelEncoder": "lowercase",
		  "timeKey": "datetime",
		  "timeEncoder": "iso8601"
		}
	  }`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(fmt.Errorf("failed to unmarshal logger config: %w", err))
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(fmt.Errorf("failed to build logger: %w", err))
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	return sugar
}
