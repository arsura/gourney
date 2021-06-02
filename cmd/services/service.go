package service

import (
	logger "github.com/arsura/moonbase-service/pkg/logger"

	"github.com/arsura/moonbase-service/pkg/models/pgsql"
)

type Service struct {
	Logger *logger.Logger
	PgRepo *pgsql.Repositories
}

type Services struct {
	Currency CurrencyService
}
