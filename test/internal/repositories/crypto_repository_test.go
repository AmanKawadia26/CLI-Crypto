package repositories_test

import (
	"cryptotracker/internal/repositories"
	"cryptotracker/pkg/config"
	"testing"

	"cryptotracker/models"
	"github.com/stretchr/testify/assert"
)

func TestPostgresCryptoRepository(t *testing.T) {
	_ = config.LoadConfig()
	conn, cleanup, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer cleanup()

	repo := repositories.NewPostgresCryptoRepository(conn)

	t.Run("DisplayTopCryptocurrencies", func(t *testing.T) {
		topCryptos, err := repo.DisplayTopCryptocurrencies()
		assert.NoError(t, err)
		assert.NotNil(t, topCryptos)
	})

	t.Run("SearchCryptocurrency", func(t *testing.T) {
		tests := []struct {
			name     string
			symbol   string
			expected float64
			wantErr  bool
		}{
			{"Found", "BTC", 10000, false},
			{"NotFound", "XYZ", 0, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				price, name, symbol, err := repo.SearchCryptocurrency(nil, &models.User{}, tt.symbol)
				if (err != nil) != tt.wantErr {
					t.Errorf("SearchCryptocurrency() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr {
					assert.NotZero(t, price, "Price should be non-zero")
					assert.NotEmpty(t, name, "Name should not be empty")
					assert.Equal(t, tt.symbol, symbol, "Symbol should match")
				}
			})
		}
	})

	t.Run("SetPriceAlert", func(t *testing.T) {
		tests := []struct {
			name        string
			symbol      string
			targetPrice float64
			wantErr     bool
		}{
			//{"Success", "BTC", 10000, false},
			{"InvalidSymbol", "XYZ", 10000, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				currentPrice, err := repo.SetPriceAlert(nil, &models.User{}, tt.symbol, tt.targetPrice)
				if (err != nil) != tt.wantErr {
					t.Errorf("SetPriceAlert() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr {
					assert.NotZero(t, currentPrice, "Current price should be non-zero")
				}
			})
		}
	})
}

func TestCheckCryptocurrencyExists_RealAPI(t *testing.T) {
	// Test case 1: Valid cryptocurrency (BTC)
	//exists, err := repositories.CheckCryptocurrencyExists("BTC")
	//assert.NoError(t, err, "error should be nil for existing cryptocurrency")
	//assert.True(t, exists, "BTC should exist")

	// Test case 2: Invalid cryptocurrency symbol
	exists, err := repositories.CheckCryptocurrencyExists("INVALID_SYMBOL")
	assert.NoError(t, err, "error should be nil for non-existing cryptocurrency")
	assert.False(t, exists, "INVALID_SYMBOL should not exist")

	// Test case 3: Empty symbol
	exists, err = repositories.CheckCryptocurrencyExists("")
	assert.NoError(t, err, "error should be nil for empty symbol")
	assert.False(t, exists, "Empty symbol should not exist")
}
