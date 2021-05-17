package fakedatabase

type ExchangeRate struct {
	Source string
	Target string
	Rate   float64
}

var ExchangeRates = []ExchangeRate{
	{
		Source: "THBT",
		Target: "MOON",
		Rate:   50,
	},
}

func FindExchangeRate(sourceCurrency string, targetCurrency string) float64 {
	var rate float64 = 0
	for _, exchangeRate := range ExchangeRates {
		if exchangeRate.Source == sourceCurrency && exchangeRate.Target == targetCurrency {
			rate = exchangeRate.Rate
			break
		}
		if exchangeRate.Source == targetCurrency && exchangeRate.Target == sourceCurrency {
			rate = 1 / exchangeRate.Rate
			break
		}
	}
	return rate
}

func UpdateExchangeRate(sourceCurrency string, targetCurrency string, rate float64) ExchangeRate {
	var newExchangeRateDate ExchangeRate
	for i, exchangeRate := range ExchangeRates {
		if exchangeRate.Source == sourceCurrency && exchangeRate.Target == targetCurrency {
			ExchangeRates[i].Rate = rate
			newExchangeRateDate = ExchangeRates[i]
			break
		}
	}
	return newExchangeRateDate
}
