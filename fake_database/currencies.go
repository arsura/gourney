package fakedatabase

type Currency struct {
	Name       string
	Amount     float64
	Total      float64
	RiseRate   float64
	RiseFactor float64
}

var Currencies = []Currency{
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

func FindCurrencyByName(name string) Currency {
	var currencyData Currency
	for _, currency := range Currencies {
		if currency.Name == name {
			currencyData = currency
			break
		}
	}
	return currencyData
}

func DecreaseCurrencyAmount(name string, decreaseVal float64) Currency {
	var currencyData Currency
	for i, currency := range Currencies {
		if currency.Name == name {
			Currencies[i].Amount -= decreaseVal
			currencyData = Currencies[i]
			break
		}
	}
	return currencyData
}
