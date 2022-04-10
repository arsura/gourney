package logger

import (
	"fmt"

	"go.uber.org/zap"
)

func NewLogger() *zap.SugaredLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Errorf("failed to build logger: %w", err))
	}
	return logger.Sugar()
}
