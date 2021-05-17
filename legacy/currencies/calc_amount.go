package currencies

import (
	"math"

	fakedatabase "github.com/arsura/moonbase-service/fake_database"
	"github.com/gofiber/fiber/v2"
)

type CalcAmountQuery struct {
	Source string  `query:"source"`
	Target string  `query:"target"`
	Amount float64 `query:"amount"`
	Type   string  `query:"type"`
}

func CalcAmountsHandler(c *fiber.Ctx) error {
	query := new(CalcAmountQuery)
	if err := c.QueryParser(query); err != nil {
		return err
	}
	sourceCurrency := query.Source
	targetCurrency := query.Target
	amountVal := query.Amount
	amountType := query.Type

	var amountResult float64 = 0
	if amountType == "target" {
		amountResult = CalcTargetAmount(sourceCurrency, targetCurrency, amountVal)
	} else if amountType == "source" {
		amountResult = CalcSourceAmount(sourceCurrency, targetCurrency, amountVal)
	}
	return c.JSON(&fiber.Map{"data": (amountResult)})
}

func CalcSourceAmount(sourceCurrency string, targetCurrency string, targetAmount float64) float64 {
	rate := fakedatabase.FindExchangeRate(sourceCurrency, targetCurrency)
	targetCurrencyData := fakedatabase.FindCurrencyByName(targetCurrency)

	var sourceAmount float64 = 0
	riseRate := targetCurrencyData.RiseRate
	riseFactor := targetCurrencyData.RiseFactor
	currentTargetAmount := targetCurrencyData.Amount
	remainder := math.Mod(currentTargetAmount, riseFactor)

	if (remainder > 0) && ((targetAmount - remainder) > 0) {
		targetAmount = targetAmount - remainder
		sourceAmount = sourceAmount + (remainder * rate)
		rate = rate + (rate * riseRate)
	}

	for targetAmount > 0 {
		if targetAmount < riseFactor {
			sourceAmount = sourceAmount + (targetAmount * rate)
			break
		}
		targetAmount = targetAmount - riseFactor
		sourceAmount = sourceAmount + (riseFactor * rate)
		rate = rate + (rate * riseRate)
	}

	return sourceAmount
}

func CalcTargetAmount(sourceCurrency string, targetCurrency string, sourceAmount float64) float64 {
	rate := fakedatabase.FindExchangeRate(sourceCurrency, targetCurrency)
	targetCurrencyData := fakedatabase.FindCurrencyByName(targetCurrency)

	var targetAmount float64 = 0
	permanentSourceAmount := sourceAmount
	riseRate := targetCurrencyData.RiseRate
	riseFactor := targetCurrencyData.RiseFactor
	currentTargetAmount := targetCurrencyData.Amount
	remainder := math.Mod(currentTargetAmount, riseFactor)

	if remainder > 0 {
		sourceAmount = sourceAmount - (rate * remainder)
		if sourceAmount > 0 {
			targetAmount = remainder
		} else {
			targetAmount = (permanentSourceAmount / rate)
		}
		rate = rate + (rate * riseRate)
	}

	for sourceAmount >= 0 {
		if sourceAmount >= (rate * riseFactor) {
			targetAmount = targetAmount + riseFactor
			sourceAmount = sourceAmount - (rate * riseFactor)
			rate = rate + (rate * riseRate)
		} else {
			targetAmount = targetAmount + (sourceAmount / rate)
			sourceAmount = sourceAmount - (rate * riseFactor)
		}
	}

	return targetAmount
}
