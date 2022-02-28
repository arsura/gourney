package api

import (
	"strconv"

	usecase "github.com/arsura/gourney/cmd/usecases"
	model "github.com/arsura/gourney/pkg/models/pgsql"
	util "github.com/arsura/gourney/pkg/utils"
	validator "github.com/arsura/gourney/pkg/validator"
	"github.com/gofiber/fiber/v2"
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
	CurrencyUsecase usecase.CurrencyUsecaseProvider
	Validator       *validator.Validator
}

func NewCurrencyHandler(currencyUsecase usecase.CurrencyUsecaseProvider, validator *validator.Validator) *currencyHandler {
	return &currencyHandler{
		CurrencyUsecase: currencyUsecase,
		Validator:       validator,
	}
}

func (h *currencyHandler) CreateCurrencyHandler(c *fiber.Ctx) error {
	currency := new(CreateReq)
	if err := c.BodyParser(currency); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": util.UnmarshalErrorParser(err),
		})
	}

	if err := h.Validator.Validate.Struct(*currency); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"errors": h.Validator.TransError(err),
		})
	}

	_, err := h.CurrencyUsecase.Create(&model.Currency{
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
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "id must be a number",
		})
	}

	result, err := h.CurrencyUsecase.FindOneById(int64(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"error": "currency not found",
		})

	}
	return c.JSON(result)
}
