package utils

import (
	"math/rand"
	"time"
)

func GenerateRandomPrices(days int, currentPrice float64) []float64 {
	rand.Seed(time.Now().UnixNano())
	prices := make([]float64, days)
	prices[days-1] = currentPrice

	for i := days - 2; i >= 0; i-- {
		// Generate a random percentage change between -5% and 5%
		change := (rand.Float64() - 0.5) * 0.1
		prices[i] = prices[i+1] * (1 + change)
	}

	return prices
}
