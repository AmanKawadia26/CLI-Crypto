package handlers

import (
	"bytes"
	"context"
	Handlers "cryptotracker/REST-API/handlers"
	mock_services "cryptotracker/test/internal/mocks/services"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test DisplayTopCryptos
func TestDisplayTopCryptos_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCryptoService := mock_services.NewMockCryptoService(ctrl)

	handler := Handlers.NewCryptoHandler(mockCryptoService)

	mockCryptoService.EXPECT().DisplayTopCryptocurrencies(10).Return([]interface{}{"BTC", "ETH"}, nil)

	req := httptest.NewRequest("GET", "/cryptos?count=10", nil)
	w := httptest.NewRecorder()

	handler.DisplayTopCryptos(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Top cryptocurrencies retrieved successfully")
}

func TestDisplayTopCryptos_InvalidCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCryptoService := mock_services.NewMockCryptoService(ctrl)

	handler := Handlers.NewCryptoHandler(mockCryptoService)

	req := httptest.NewRequest("GET", "/cryptos?count=abc", nil)
	w := httptest.NewRecorder()

	handler.DisplayTopCryptos(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request payload")
}

func TestDisplayTopCryptos_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCryptoService := mock_services.NewMockCryptoService(ctrl)

	handler := Handlers.NewCryptoHandler(mockCryptoService)

	mockCryptoService.EXPECT().DisplayTopCryptocurrencies(10).Return(nil, errors.New("service error"))

	req := httptest.NewRequest("GET", "/cryptos?count=10", nil)
	w := httptest.NewRecorder()

	handler.DisplayTopCryptos(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error retrieving top cryptocurrencies")
}

// Test DisplayCryptoByName
func TestDisplayCryptoByName_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCryptoService := mock_services.NewMockCryptoService(ctrl)

	handler := Handlers.NewCryptoHandler(mockCryptoService)

	mockCryptoService.EXPECT().
		SearchCryptocurrency(gomock.Any(), "BTC").
		Return(50000.00, "Bitcoin", "BTC", nil)

	req := httptest.NewRequest("GET", "/cryptos/BTC", nil)
	req = req.WithContext(setContext(req.Context(), "username", "testuser"))
	w := httptest.NewRecorder()

	handler.DisplayCryptoByName(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Cryptocurrency details retrieved successfully")
}

func TestDisplayCryptoByName_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCryptoService := mock_services.NewMockCryptoService(ctrl)

	handler := Handlers.NewCryptoHandler(mockCryptoService)

	mockCryptoService.EXPECT().
		SearchCryptocurrency(gomock.Any(), "BTC").
		Return(0.0, "", "", errors.New("not found"))

	req := httptest.NewRequest("GET", "/cryptos/BTC", nil)
	req = req.WithContext(setContext(req.Context(), "username", "testuser"))
	w := httptest.NewRecorder()

	handler.DisplayCryptoByName(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error searching for cryptocurrency")
}

func TestDisplayCryptoByName_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCryptoService := mock_services.NewMockCryptoService(ctrl)

	handler := Handlers.NewCryptoHandler(mockCryptoService)

	req := httptest.NewRequest("GET", "/cryptos/BTC", nil)
	w := httptest.NewRecorder()

	handler.DisplayCryptoByName(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Unauthorized access")
}

type AlertReq struct {
	Symbol      string  `json:"crypto_symbol"`
	TargetPrice float64 `json:"target_price"`
}

// Test SetPriceAlert
func TestSetPriceAlert_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCryptoService := mock_services.NewMockCryptoService(ctrl)

	handler := Handlers.NewCryptoHandler(mockCryptoService)

	alertReq := AlertReq{Symbol: "BTC", TargetPrice: 60000.00}
	mockCryptoService.EXPECT().
		SetPriceAlert(gomock.Any(), "BTC", 60000.00).
		Return(50000.00, nil)

	body, _ := json.Marshal(alertReq)
	req := httptest.NewRequest("POST", "/cryptos/alert", bytes.NewBuffer(body))
	req = req.WithContext(setContext(req.Context(), "username", "testuser"))
	w := httptest.NewRecorder()

	handler.SetPriceAlert(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Price alert set successfully")
}

func TestSetPriceAlert_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCryptoService := mock_services.NewMockCryptoService(ctrl)

	handler := Handlers.NewCryptoHandler(mockCryptoService)

	alertReq := AlertReq{Symbol: "BTC", TargetPrice: 60000.00}
	mockCryptoService.EXPECT().
		SetPriceAlert(gomock.Any(), "BTC", 60000.00).
		Return(0.0, errors.New("service error"))

	body, _ := json.Marshal(alertReq)
	req := httptest.NewRequest("POST", "/cryptos/alert", bytes.NewBuffer(body))
	req = req.WithContext(setContext(req.Context(), "username", "testuser"))
	w := httptest.NewRecorder()

	handler.SetPriceAlert(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error setting price alert")
}

// Helper function to set context
func setContext(ctx context.Context, key, value string) context.Context {
	return context.WithValue(ctx, key, value)
}
