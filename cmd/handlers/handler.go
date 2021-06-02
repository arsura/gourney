package handler

import (
	service "github.com/arsura/moonbase-service/cmd/services"
	validator "github.com/arsura/moonbase-service/pkg/validator"
)

type Handler struct {
	Services  *service.Services
	Validator *validator.Validator
}
