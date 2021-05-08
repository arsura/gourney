package currencies

import (
	"math"

	fakedatabase "github.com/arsura/moonbase-service/fake_database"
	"github.com/gofiber/fiber/v2"
)

type PurchaseReqBody struct {
	Source            string  `json:"source"`
	Target            string  `json:"target"`
	SourceAmount      float64 `json:"sourceAmount"`
	TargetAmount      float64 `json:"targetAmount"`
	SlippageTolerance float64 `json:"slippage"`
	ClientId          string  `json:"clientId"`
}

func PurchaseHandler(c *fiber.Ctx) error {
	purchaseData := new(PurchaseReqBody)
	if err := c.BodyParser(purchaseData); err != nil {
		return err
	}
	sourceCurrency := purchaseData.Source
	targetCurrency := purchaseData.Target
	sourceAmount := purchaseData.SourceAmount
	expectedTargetAmount := purchaseData.TargetAmount
	slippageTolerance := purchaseData.SlippageTolerance
	clientId := purchaseData.ClientId
	lowerBound, upperBound := CalcToleranceBound(slippageTolerance, expectedTargetAmount)
	actualTargetAmount := CalcTargetAmount(sourceCurrency, targetCurrency, sourceAmount)
	if actualTargetAmount < lowerBound || actualTargetAmount > upperBound {
		return fiber.ErrBadRequest
	}
	newRate := AdjustCurrency(sourceCurrency, targetCurrency, actualTargetAmount)
	fakedatabase.UpdateHistoryLog(clientId, sourceCurrency, sourceAmount, targetCurrency, actualTargetAmount, newRate)
	return c.JSON(&fiber.Map{"data": actualTargetAmount})
}

func CalcToleranceBound(slippageTolerance float64, expectedTargetAmount float64) (float64, float64) {
	upperBound := expectedTargetAmount + (expectedTargetAmount * slippageTolerance)
	lowerBound := expectedTargetAmount - (expectedTargetAmount * slippageTolerance)
	return lowerBound, upperBound
}

func CalcNewCurencyRate(currentRate float64, prevCurrencyData fakedatabase.Currency, newCurrencyData fakedatabase.Currency) float64 {
	var newRate float64 = currentRate
	riseTimes := int(((math.Ceil(prevCurrencyData.Amount/10) * 10) - (math.Ceil(newCurrencyData.Amount/10) * 10)) / newCurrencyData.RiseFactor)
	for i := 0; i < riseTimes; i++ {
		newRate = newRate + (newRate * newCurrencyData.RiseRate)
	}
	return newRate
}

func AdjustCurrency(sourceCurrency string, targetCurrency string, decreaseAmount float64) float64 {
	prevCurrencyData := fakedatabase.FindCurrencyByName(targetCurrency)
	currentRate := fakedatabase.FindExchangeRate(sourceCurrency, targetCurrency)
	newCurrencyData := fakedatabase.DecreaseCurrencyAmount(targetCurrency, decreaseAmount)
	newRate := CalcNewCurencyRate(currentRate, prevCurrencyData, newCurrencyData)
	fakedatabase.UpdateExchangeRate(sourceCurrency, targetCurrency, newRate)
	return newRate
}
