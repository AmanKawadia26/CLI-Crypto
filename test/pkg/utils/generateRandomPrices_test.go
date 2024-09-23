package utils

import (
	"cryptotracker/pkg/utils"
	"math"
	"testing"
)

func TestGenerateRandomPrices(t *testing.T) {
	tests := []struct {
		name         string
		days         int
		currentPrice float64
	}{
		{
			name:         "Single day",
			days:         1,
			currentPrice: 100.0,
		},
		{
			name:         "Two days",
			days:         2,
			currentPrice: 100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prices := utils.GenerateRandomPrices(tt.days, tt.currentPrice)

			// Check if the length of prices slice is correct
			if len(prices) != tt.days {
				t.Errorf("GenerateRandomPrices() length = %v, want %v", len(prices), tt.days)
			}

			// Check if the last price is the current price
			if math.Abs(prices[len(prices)-1]-tt.currentPrice) > 1e-6 {
				t.Errorf("GenerateRandomPrices() last price = %v, want %v", prices[len(prices)-1], tt.currentPrice)
			}

			// Check if the prices are valid floats
			for _, price := range prices {
				if math.IsNaN(price) || math.IsInf(price, 0) {
					t.Errorf("GenerateRandomPrices() returned invalid price: %v", price)
				}
			}

			// Check if the prices decrease or increase by a reasonable amount
			for i := 0; i < len(prices)-1; i++ {
				change := (prices[i+1] - prices[i]) / prices[i]
				if math.Abs(change) > 0.05 {
					t.Errorf("GenerateRandomPrices() percentage change = %v, is greater than expected", change)
				}
			}
		})
	}
}
