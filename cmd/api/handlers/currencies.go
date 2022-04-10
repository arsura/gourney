package api

import (
	"context"
	"encoding/json"
	"strconv"

	usecase "github.com/arsura/gourney/cmd/usecases"
	"github.com/arsura/gourney/pkg/constant"
	model "github.com/arsura/gourney/pkg/models/pgsql"
	util "github.com/arsura/gourney/pkg/utils"
	validator "github.com/arsura/gourney/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type CreateReq struct {
	Name       string  `json:"name" validate:"required"`
	Amount     float64 `json:"amount" validate:"required"`
	Total      float64 `json:"total" validate:"required"`
	RiseRate   float64 `json:"rise_rate"`
	RiseFactor float64 `json:"rise_factor"`
}

type CurrencyHandlerProvider interface {
	CreateCurrencyHandler(c *fiber.Ctx) error
	FindCurrencyByIdHandler(c *fiber.Ctx) error
}

type currencyHandler struct {
	currencyUsecase usecase.CurrencyUsecaseProvider
	validator       *validator.Validator
	logger          *zap.SugaredLogger
}

func NewCurrencyHandler(currencyUsecase usecase.CurrencyUsecaseProvider, validator *validator.Validator, logger *zap.SugaredLogger) *currencyHandler {
	return &currencyHandler{
		currencyUsecase: currencyUsecase,
		validator:       validator,
		logger:          logger,
	}
}

func (h *currencyHandler) CreateCurrencyHandler(c *fiber.Ctx) error {
	var (
		currency  = &CreateReq{}
		requestId = c.Locals("requestid").(string)
		ctx       = context.WithValue(c.UserContext(), constant.RequestIdKey, requestId)
	)

	err := json.Unmarshal(c.Body(), currency)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": util.UnmarshalErrorParser(err),
		})
	}

	if err := h.validator.Validate.Struct(*currency); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": h.validator.TransError(err),
		})
	}

	_, err = h.currencyUsecase.Create(ctx, &model.Currency{
		Name:       currency.Name,
		Amount:     currency.Amount,
		Total:      currency.Total,
		RiseRate:   currency.RiseRate,
		RiseFactor: currency.RiseFactor,
	})

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "failed to create currency",
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *currencyHandler) FindCurrencyByIdHandler(c *fiber.Ctx) error {
	var (
		requestId = c.Locals("requestid").(string)
		ctx       = context.WithValue(c.UserContext(), constant.RequestIdKey, requestId)
	)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "id must be a number",
		})
	}

	result, err := h.currencyUsecase.FindOneById(ctx, int64(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"error": "currency not found",
		})

	}
	return c.JSON(result)
}
