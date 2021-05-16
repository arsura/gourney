package models

import (
	"errors"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Currency struct {
	ID         int
	Name       string
	Amount     float64
	Total      float64
	RiseRate   float64
	RiseFactor float64
}
