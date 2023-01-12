package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DATA_FIELD        = "data"
	ERROR_FIELD       = "error"
	TRACKING_ID_FIELD = "tracking_id"
)

func NewLogger() *zap.SugaredLogger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.FunctionKey = "func"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		panic(fmt.Errorf("failed to new logger: %w", err))
	}
	return logger.Sugar()
}
