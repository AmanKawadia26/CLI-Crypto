//go:build !test
// +build !test

package repositories

import (
	"cryptotracker/internal/api"
	"cryptotracker/models"
	"cryptotracker/pkg/storage"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/jackc/pgx/v4"
	//"github.com/jackc/pgx/v5"
	"strings"
	"time"
)

type CryptoRepository interface {
	DisplayTopCryptocurrencies() ([]interface{}, error)
	SearchCryptocurrency(conn *pgx.Conn, user *models.User, cryptoSymbol string) (float64, string, string, error)
	SetPriceAlert(conn *pgx.Conn, user *models.User, symbol string, targetPrice float64) (float64, error)
}

type PostgresCryptoRepository struct {
	conn *pgx.Conn
}

func NewPostgresCryptoRepository(conn *pgx.Conn) *PostgresCryptoRepository {
	return &PostgresCryptoRepository{
		conn: conn,
	}
}

func (repo *PostgresCryptoRepository) DisplayTopCryptocurrencies() ([]interface{}, error) {
	params := map[string]string{
		"start":   "1",
		"limit":   "10",
		"convert": "USD",
	}

	response := api.GetAPIResponse("/listings/latest", params)

	//fmt.Println(response)

	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("error unmarshalling API response: %v", err)
	}

	//fmt.Println(result)

	data, ok := result["data"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("data not found in the response")
	}

	return data, nil
}

func (repo *PostgresCryptoRepository) SearchCryptocurrency(conn *pgx.Conn, user *models.User, cryptoSymbol string) (float64, string, string, error) {

	cryptoSymbol = strings.ToLower(cryptoSymbol)

	// Define parameters for the API request
	params := map[string]string{
		"start":   "1",
		"limit":   "5000",
		"convert": "USD",
	}

	// Make the API request
	response := api.GetAPIResponse("/listings/latest", params)

	// Parse the API response
	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		return 0, "", "", fmt.Errorf("error unmarshalling API response: %v", err)
	}

	// Validate if 'data' field exists in the response
	data, ok := result["data"].([]interface{})
	if !ok || len(data) == 0 {
		return 0, "", "", fmt.Errorf("data not found or empty in the response")
	}

	// Loop through the cryptocurrencies to find a match
	for _, item := range data {
		crypto, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		symbol, ok := crypto["symbol"].(string)
		if !ok {
			continue
		}

		name, ok := crypto["name"].(string)
		if !ok {
			continue
		}

		// Check if the symbol or name matches the input
		if strings.ToLower(symbol) == cryptoSymbol || strings.ToLower(name) == cryptoSymbol {
			// Verify if the price field exists in the response
			quote, ok := crypto["quote"].(map[string]interface{})
			if !ok {
				return 0, "", "", fmt.Errorf("quote data not found for cryptocurrency")
			}

			usd, ok := quote["USD"].(map[string]interface{})
			if !ok {
				return 0, "", "", fmt.Errorf("USD data not found in the quote")
			}

			price, ok := usd["price"].(float64)
			if !ok {
				return 0, "", "", fmt.Errorf("could not retrieve or convert the price")
			}

			// Return the cryptocurrency's price, name, and symbol if found
			return price, name, symbol, nil
		}
	}

	// Log and inform the user if the cryptocurrency was not found
	color.New(color.FgYellow).Printf("Cryptocurrency not found for input: %s\n", cryptoSymbol)
	color.New(color.FgMagenta).Println("Please request the addition of this cryptocurrency to our app.")

	// Save a request to add the unavailable cryptocurrency
	request := &models.UnavailableCryptoRequest{
		CryptoSymbol:   cryptoSymbol,
		UserName:       user.Username,
		RequestMessage: "Please add this cryptocurrency.",
		Status:         "Pending",
		Timestamp:      time.Now(),
	}

	// Ensure the database connection is valid before attempting to save
	if conn == nil {
		return 0, "", "", fmt.Errorf("database connection is nil")
	}

	unavailableCryptoRepo := storage.NewPGUnavailableCryptoRequestRepository(repo.conn)

	// Save the unavailable crypto request in the database
	if err := unavailableCryptoRepo.SaveUnavailableCryptoRequest(conn, request); err != nil {
		return 0, "", "", fmt.Errorf("error saving unavailable crypto request: %v", err)
	}

	// Return a message indicating the cryptocurrency was not found but a request was made
	return 0, "", "", fmt.Errorf("request to add the cryptocurrency has been submitted")
}

func (repo *PostgresCryptoRepository) SetPriceAlert(conn *pgx.Conn, user *models.User, symbol string, targetPrice float64) (float64, error) {
	// Step 1: Check if the cryptocurrency exists by calling the API
	exists, err := CheckCryptocurrencyExists(symbol)
	if err != nil {
		return 0, fmt.Errorf("error checking cryptocurrency existence: %v", err)
	}

	// If the cryptocurrency does not exist, return an error
	if !exists {
		return 0, fmt.Errorf("cryptocurrency with symbol %s does not exist", symbol)
	}

	// Step 2: Continue with the original price alert logic
	params := map[string]string{
		"symbol":  symbol,
		"convert": "USD",
	}

	response := api.GetAPIResponse("/quotes/latest", params)
	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		return 0, fmt.Errorf("error unmarshalling API response: %v", err)
	}

	data, dataOk := result["data"].(map[string]interface{})
	if !dataOk || data[symbol] == nil {
		return 0, fmt.Errorf("cryptocurrency data not found for symbol: %s", symbol)
	}

	cryptoData, ok := data[symbol].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("unexpected data structure for symbol: %s", symbol)
	}

	quote, ok := cryptoData["quote"].(map[string]interface{})
	if !ok || quote["USD"] == nil {
		return 0, fmt.Errorf("quote data not available for %s", symbol)
	}

	priceData, ok := quote["USD"].(map[string]interface{})
	if !ok || priceData["price"] == nil {
		return 0, fmt.Errorf("price data not available for %s", symbol)
	}

	cryptoIDInterface, ok := cryptoData["id"].(float64)
	if !ok || cryptoIDInterface == 0 {
		return 0, fmt.Errorf("cryptocurrency ID not found for symbol: %s", symbol)
	}

	cryptoID := int(cryptoIDInterface)
	currentPrice, ok := priceData["price"].(float64)
	if !ok {
		return 0, fmt.Errorf("failed to convert price to float64")
	}

	if currentPrice >= targetPrice {
		return currentPrice, fmt.Errorf("alert: %s has reached your target price of $%.2f. Current price: $%.2f", symbol, targetPrice, currentPrice)
	}

	// Save the notification in the database
	notification := &models.PriceNotification{
		CryptoID:    cryptoID,
		Crypto:      symbol,
		TargetPrice: targetPrice,
		Username:    user.Username,
		AskedAt:     time.Now().Format(time.RFC3339),
		Status:      "Pending",
	}

	notificationRepo := storage.NewPGNotificationRepository(repo.conn)
	if err := notificationRepo.SavePriceNotification(conn, notification); err != nil {
		return 0, fmt.Errorf("error saving notification: %v", err)
	}

	return currentPrice, fmt.Errorf("%s is still below your target price. Current price: $%.2f. Notification created.", symbol, currentPrice)
}

func CheckCryptocurrencyExists(symbol string) (bool, error) {
	params := map[string]string{
		"symbol":  symbol,
		"convert": "USD",
	}

	response := api.GetAPIResponse("/quotes/latest", params)
	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		return false, fmt.Errorf("error unmarshalling API response: %v", err)
	}

	// Check if the data contains the requested symbol
	data, dataOk := result["data"].(map[string]interface{})
	if !dataOk || data[symbol] == nil {
		return false, nil // Cryptocurrency does not exist
	}

	return true, nil // Cryptocurrency exists
}
