package main

import (
	"log"
	"math"
	"time"

	"github.com/gofiber/fiber/v2"
)

type CalcAmountQuery struct {
	Source string  `query:"source"`
	Target string  `query:"target"`
	Amount float64 `query:"amount"`
	Type   string  `query:"type"`
}

type GetExchangeRateQuery struct {
	Source string `query:"source"`
	Target string `query:"target"`
}

type GetHistoryLogsQuery struct {
	Source string `query:"source"`
	Target string `query:"target"`
	Skip   int32  `query:"skip"`
	Limit  int32  `query:"limit"`
	Sort   string `query:"sort"`
}

type PurchaseBody struct {
	Source            string  `json:"source"`
	Target            string  `json:"target"`
	SourceAmount      float64 `json:"sourceAmount"`
	TargetAmount      float64 `json:"targetAmount"`
	SlippageTolerance float64 `json:"slippage"`
	ClientId          string  `json:"clientId"`
}

type ExchangeRate struct {
	Source string
	Target string
	Rate   float64
}

type Currency struct {
	Name       string
	Amount     float64
	Total      float64
	RiseRate   float64
	RiseFactor float64
}

type HistoryLog struct {
	ClientId     string
	Source       string
	Target       string
	SourceAmount float64
	TargetAmount float64
	Rate         float64
	CreatedAt    time.Time
}

var exchangeRates []ExchangeRate
var currencies []Currency
var historyLogs []HistoryLog

func updateHistoryLog(clientId string, source string, sourceAmount float64, target string, targetAmount float64, currentRate float64) {
	historyLogs =
		append(
			historyLogs,
			HistoryLog{
				ClientId:     clientId,
				Source:       source,
				Target:       target,
				SourceAmount: sourceAmount,
				TargetAmount: targetAmount,
				Rate:         currentRate,
				CreatedAt:    time.Now().UTC(),
			},
		)
}

func calcTargetAmount(sourceCurrency string, targetCurrency string, sourceAmount float64) (float64, float64) {
	var rate float64 = 0
	for _, exchangeRate := range exchangeRates {
		if exchangeRate.Source == sourceCurrency && exchangeRate.Target == targetCurrency {
			rate = exchangeRate.Rate
			break
		}
	}

	var targetCurrencyData Currency

	for _, currency := range currencies {
		if currency.Name == targetCurrency {
			targetCurrencyData = currency
			break
		}
	}

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

	return targetAmount, rate
}

func calcSourceAmount(sourceCurrency string, targetCurrency string, targetAmount float64) (float64, float64) {
	var rate float64 = 0
	for _, exchangeRate := range exchangeRates {
		if exchangeRate.Source == sourceCurrency && exchangeRate.Target == targetCurrency {
			rate = exchangeRate.Rate
			break
		}
	}

	var targetCurrencyData Currency

	for _, currency := range currencies {
		if currency.Name == targetCurrency {
			targetCurrencyData = currency
			break
		}
	}

	var sourceAmount float64 = 0
	riseRate := targetCurrencyData.RiseRate
	riseFactor := targetCurrencyData.RiseFactor
	currentTargetAmount := targetCurrencyData.Amount
	remainder := math.Mod(currentTargetAmount, riseFactor)

	if remainder > 0 {
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

	return sourceAmount, rate
}

func adjustCurrency(sourceCurrency string, targetCurrency string, decreaseAmount float64, newRate float64) {
	for i, exchangeRate := range exchangeRates {
		if exchangeRate.Source == sourceCurrency && exchangeRate.Target == targetCurrency {
			exchangeRates[i].Rate = newRate
			break
		}
	}

	for i, currency := range currencies {
		if currency.Name == targetCurrency {
			currencies[i].Amount -= decreaseAmount
			break
		}
	}
	return
}

func calcToleranceBound(slippageTolerance float64, expectedTargetAmount float64) (float64, float64) {
	upperBound := expectedTargetAmount + (expectedTargetAmount * slippageTolerance)
	lowerBound := expectedTargetAmount - (expectedTargetAmount * slippageTolerance)
	return lowerBound, upperBound
}

func purchaseHandler(c *fiber.Ctx) error {
	purchaseData := new(PurchaseBody)
	if err := c.BodyParser(purchaseData); err != nil {
		return err
	}
	sourceCurrency := purchaseData.Source
	targetCurrency := purchaseData.Target
	sourceAmount := purchaseData.SourceAmount
	expectedTargetAmount := purchaseData.TargetAmount
	slippageTolerance := purchaseData.SlippageTolerance
	clientId := purchaseData.ClientId
	lowerBound, upperBound := calcToleranceBound(slippageTolerance, expectedTargetAmount)
	actualTargetAmount, newRate := calcTargetAmount(sourceCurrency, targetCurrency, sourceAmount)
	if actualTargetAmount < lowerBound || actualTargetAmount > upperBound {
		return fiber.ErrBadRequest
	}
	adjustCurrency(sourceCurrency, targetCurrency, actualTargetAmount, newRate)
	updateHistoryLog(clientId, sourceCurrency, sourceAmount, targetCurrency, actualTargetAmount, newRate)
	return c.JSON(&fiber.Map{"data": actualTargetAmount})
}

func calcAmountHandler(c *fiber.Ctx) error {
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
		amountResult, _ = calcTargetAmount(sourceCurrency, targetCurrency, amountVal)
	} else if amountType == "source" {
		amountResult, _ = calcSourceAmount(sourceCurrency, targetCurrency, amountVal)
	}
	return c.JSON(&fiber.Map{"data": (amountResult)})
}

func getExchangeRateHandler(c *fiber.Ctx) error {
	var rate float64 = 0
	query := new(CalcAmountQuery)
	if err := c.QueryParser(query); err != nil {
		return err
	}
	sourceCurrency := query.Source
	targetCurrency := query.Target
	for _, exchangeRate := range exchangeRates {
		if exchangeRate.Source == sourceCurrency && exchangeRate.Target == targetCurrency {
			rate = exchangeRate.Rate
		}
		if exchangeRate.Source == targetCurrency && exchangeRate.Target == sourceCurrency {
			rate = 1 / exchangeRate.Rate
		}
	}
	return c.JSON(&fiber.Map{"data": rate})
}

func getHistoryLogsHandler(c *fiber.Ctx) error {
	query := new(GetHistoryLogsQuery)
	if err := c.QueryParser(query); err != nil {
		return err
	}
	// sourceCurrency := query.Source
	// targetCurrency := query.Target
	// skip := query.Skip
	// limit := query.Limit
	// sort := query.Sort
	return c.JSON(&fiber.Map{"data": []HistoryLog(historyLogs)})
}

func main() {
	currencies = []Currency{
		{
			Name:       "MOON",
			Amount:     1000,
			Total:      1000,
			RiseRate:   0.1,
			RiseFactor: 10,
		},
		{
			Name:       "THBT",
			Amount:     0,
			Total:      0,
			RiseRate:   0,
			RiseFactor: 0,
		},
	}

	exchangeRates = []ExchangeRate{
		{
			Source: "THBT",
			Target: "MOON",
			Rate:   50,
		},
		{
			Source: "THBT",
			Target: "THB",
			Rate:   1,
		},
	}

	historyLogs = []HistoryLog{}
	app := fiber.New()

	app.Get("/currencies/:name", func(c *fiber.Ctx) error {
		currencyName := c.Params("name")
		var currencyData Currency
		for _, currency := range currencies {
			if currency.Name == currencyName {
				currencyData = currency
				break
			}
		}
		return c.JSON(&fiber.Map{"data": Currency(currencyData)})
	})
	app.Get("/amounts", calcAmountHandler)
	app.Get("/rates", getExchangeRateHandler)
	app.Get("/logs", getHistoryLogsHandler)
	app.Post("/purchases", purchaseHandler)

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
