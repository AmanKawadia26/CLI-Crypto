package services_test

import (
	"cryptotracker/internal/services"
	"cryptotracker/models"
	mock_repositories "cryptotracker/test/internal/mocks/repository"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	"reflect"
	"testing"
)

func TestNewCryptoService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repositories.NewMockCryptoRepository(ctrl)
	service := services.NewCryptoService(mockRepo)

	if service == nil {
		t.Error("Expected non-nil CryptoService")
	}

	if reflect.TypeOf(service) != reflect.TypeOf(&services.CryptoServiceImpl{}) {
		t.Error("Expected CryptoServiceImpl type")
	}
}

func TestCryptoServiceImpl_DisplayTopCryptocurrencies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repositories.NewMockCryptoRepository(ctrl)
	service := services.NewCryptoService(mockRepo)

	testCases := []struct {
		name          string
		mockResult    []interface{}
		mockErr       error
		expectedError bool
	}{
		{
			name:          "Success",
			mockResult:    []interface{}{map[string]interface{}{"symbol": "BTC", "price": 50000.0}},
			mockErr:       nil,
			expectedError: false,
		},
		{
			name:          "Error",
			mockResult:    nil,
			mockErr:       errors.New("API error"),
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.EXPECT().DisplayTopCryptocurrencies().Return(tc.mockResult, tc.mockErr)

			result, err := service.DisplayTopCryptocurrencies()

			if tc.expectedError && err == nil {
				t.Error("Expected an error, got nil")
			}

			if !tc.expectedError && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}

			if !reflect.DeepEqual(result, tc.mockResult) {
				t.Errorf("Expected result %v, got %v", tc.mockResult, result)
			}
		})
	}
}

func TestCryptoServiceImpl_SearchCryptocurrency(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repositories.NewMockCryptoRepository(ctrl)
	service := services.NewCryptoService(mockRepo)

	testUser := &models.User{Username: "testuser"}
	testConn := &pgx.Conn{}

	testCases := []struct {
		name          string
		cryptoSymbol  string
		mockPrice     float64
		mockName      string
		mockCoinID    string
		mockErr       error
		expectedError bool
	}{
		{
			name:          "Success",
			cryptoSymbol:  "BTC",
			mockPrice:     50000.0,
			mockName:      "Bitcoin",
			mockCoinID:    "bitcoin",
			mockErr:       nil,
			expectedError: false,
		},
		{
			name:          "Error",
			cryptoSymbol:  "INVALID",
			mockPrice:     0,
			mockName:      "",
			mockCoinID:    "",
			mockErr:       errors.New("Cryptocurrency not found"),
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.EXPECT().SearchCryptocurrency(testConn, testUser, tc.cryptoSymbol).Return(tc.mockPrice, tc.mockName, tc.mockCoinID, tc.mockErr)

			price, name, coinID, err := service.SearchCryptocurrency(testConn, testUser, tc.cryptoSymbol)

			if tc.expectedError && err == nil {
				t.Error("Expected an error, got nil")
			}

			if !tc.expectedError && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}

			if price != tc.mockPrice {
				t.Errorf("Expected price %v, got %v", tc.mockPrice, price)
			}

			if name != tc.mockName {
				t.Errorf("Expected name %v, got %v", tc.mockName, name)
			}

			if coinID != tc.mockCoinID {
				t.Errorf("Expected coinID %v, got %v", tc.mockCoinID, coinID)
			}
		})
	}
}

func TestCryptoServiceImpl_SetPriceAlert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repositories.NewMockCryptoRepository(ctrl)
	service := services.NewCryptoService(mockRepo)

	testUser := &models.User{Username: "testuser"}
	testConn := &pgx.Conn{}

	testCases := []struct {
		name             string
		symbol           string
		targetPrice      float64
		mockCurrentPrice float64
		mockErr          error
		expectedError    bool
	}{
		{
			name:             "Success",
			symbol:           "BTC",
			targetPrice:      60000.0,
			mockCurrentPrice: 50000.0,
			mockErr:          nil,
			expectedError:    false,
		},
		{
			name:             "Error",
			symbol:           "INVALID",
			targetPrice:      1000.0,
			mockCurrentPrice: 0,
			mockErr:          errors.New("Failed to set price alert"),
			expectedError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.EXPECT().SetPriceAlert(testConn, testUser, tc.symbol, tc.targetPrice).Return(tc.mockCurrentPrice, tc.mockErr)

			currentPrice, err := service.SetPriceAlert(testConn, testUser, tc.symbol, tc.targetPrice)

			if tc.expectedError && err == nil {
				t.Error("Expected an error, got nil")
			}

			if !tc.expectedError && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}

			if currentPrice != tc.mockCurrentPrice {
				t.Errorf("Expected current price %v, got %v", tc.mockCurrentPrice, currentPrice)
			}
		})
	}
}
