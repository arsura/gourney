package model

type Currency struct {
	Id         int64   `json:"id"`
	Name       string  `json:"name"`
	Amount     float64 `json:"amount"`
	Total      float64 `json:"total"`
	RiseRate   float64 `json:"rise_rate"`
	RiseFactor float64 `json:"rise_factor"`
}
