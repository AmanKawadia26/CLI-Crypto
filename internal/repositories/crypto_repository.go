//package repositories
//
//import (
//	"cryptotracker/internal/api"
//	"cryptotracker/models"
//	"cryptotracker/pkg/storage"
//	"encoding/json"
//	"fmt"
//	"github.com/fatih/color"
//	"github.com/jackc/pgx/v4"
//	//"github.com/jackc/pgx/v5"
//	"strings"
//	"time"
//)
//
//type CryptoRepository interface {
//	DisplayTopCryptocurrencies() ([]interface{}, error)
//	SearchCryptocurrency(user *models.User, cryptoSymbol string) (float64, string, string, error)
//	SetPriceAlert(user *models.User, symbol string, targetPrice float64) (float64, error)
//}
//
//type PostgresCryptoRepository struct {
//	conn *pgx.Conn
//}
//
//func NewPostgresCryptoRepository(conn *pgx.Conn) *PostgresCryptoRepository {
//	return &PostgresCryptoRepository{
//		conn: conn,
//	}
//}
//
//func (repo *PostgresCryptoRepository) DisplayTopCryptocurrencies() ([]interface{}, error) {
//	params := map[string]string{
//		"start":   "1",
//		"limit":   "10",
//		"convert": "USD",
//	}
//
//	response := api.GetAPIResponse("/listings/latest", params)
//
//	//fmt.Println(response)
//
//	var result map[string]interface{}
//	if err := json.Unmarshal(response, &result); err != nil {
//		return nil, fmt.Errorf("error unmarshalling API response: %v", err)
//	}
//
//	//fmt.Println(result)
//
//	data, ok := result["data"].([]interface{})
//	if !ok {
//		return nil, fmt.Errorf("data not found in the response")
//	}
//
//	return data, nil
//}
//
//func (repo *PostgresCryptoRepository) SearchCryptocurrency(user *models.User, cryptoSymbol string) (float64, string, string, error) {
//
//	cryptoSymbol = strings.ToLower(cryptoSymbol)
//
//	params := map[string]string{
//		"start":   "1",
//		"limit":   "5000",
//		"convert": "USD",
//	}
//
//	response := api.GetAPIResponse("/listings/latest", params)
//
//	var result map[string]interface{}
//	if err := json.Unmarshal(response, &result); err != nil {
//		return 0, "", "", fmt.Errorf("error unmarshalling API response: %v", err)
//	}
//
//	data, ok := result["data"].([]interface{})
//	if !ok || len(data) == 0 {
//		return 0, "", "", fmt.Errorf("data not found or empty in the response")
//	}
//
//	for _, item := range data {
//		crypto, ok := item.(map[string]interface{})
//		if !ok {
//			continue
//		}
//
//		symbol, ok := crypto["symbol"].(string)
//		if !ok {
//			continue
//		}
//
//		name, ok := crypto["name"].(string)
//		if !ok {
//			continue
//		}
//
//		if strings.ToLower(symbol) == cryptoSymbol || strings.ToLower(name) == cryptoSymbol {
//			quote, ok := crypto["quote"].(map[string]interface{})
//			if !ok {
//				return 0, "", "", fmt.Errorf("quote data not found for cryptocurrency")
//			}
//
//			usd, ok := quote["USD"].(map[string]interface{})
//			if !ok {
//				return 0, "", "", fmt.Errorf("USD data not found in the quote")
//			}
//
//			price, ok := usd["price"].(float64)
//			if !ok {
//				return 0, "", "", fmt.Errorf("could not retrieve or convert the price")
//			}
//
//			return price, name, symbol, nil
//		}
//	}
//
//	color.New(color.FgYellow).Printf("Cryptocurrency not found for input: %s\n", cryptoSymbol)
//	color.New(color.FgMagenta).Println("Please request the addition of this cryptocurrency to our app.")
//
//	request := &models.UnavailableCryptoRequest{
//		CryptoSymbol:   cryptoSymbol,
//		UserName:       user.Username,
//		RequestMessage: "Please add this cryptocurrency.",
//		Status:         "Pending",
//		Timestamp:      time.Now(),
//	}
//
//	if repo.conn == nil {
//		return 0, "", "", fmt.Errorf("database connection is nil")
//	}
//
//	unavailableCryptoRepo := storage.NewPGUnavailableCryptoRequestRepository(repo.conn)
//
//	if err := unavailableCryptoRepo.SaveUnavailableCryptoRequest(repo.conn, request); err != nil {
//		return 0, "", "", fmt.Errorf("error saving unavailable crypto request: %v", err)
//	}
//
//	return 0, "", "", fmt.Errorf("request to add the cryptocurrency has been submitted")
//}
//
//func (repo *PostgresCryptoRepository) SetPriceAlert(user *models.User, symbol string, targetPrice float64) (float64, error) {
//	exists, err := CheckCryptocurrencyExists(symbol)
//	if err != nil {
//		return 0, fmt.Errorf("error checking cryptocurrency existence: %v", err)
//	}
//
//	if !exists {
//		return 0, fmt.Errorf("cryptocurrency with symbol %s does not exist", symbol)
//	}
//
//	params := map[string]string{
//		"symbol":  symbol,
//		"convert": "USD",
//	}
//
//	response := api.GetAPIResponse("/quotes/latest", params)
//	var result map[string]interface{}
//	if err := json.Unmarshal(response, &result); err != nil {
//		return 0, fmt.Errorf("error unmarshalling API response: %v", err)
//	}
//
//	data, dataOk := result["data"].(map[string]interface{})
//	if !dataOk || data[symbol] == nil {
//		return 0, fmt.Errorf("cryptocurrency data not found for symbol: %s", symbol)
//	}
//
//	cryptoData, ok := data[symbol].(map[string]interface{})
//	if !ok {
//		return 0, fmt.Errorf("unexpected data structure for symbol: %s", symbol)
//	}
//
//	quote, ok := cryptoData["quote"].(map[string]interface{})
//	if !ok || quote["USD"] == nil {
//		return 0, fmt.Errorf("quote data not available for %s", symbol)
//	}
//
//	priceData, ok := quote["USD"].(map[string]interface{})
//	if !ok || priceData["price"] == nil {
//		return 0, fmt.Errorf("price data not available for %s", symbol)
//	}
//
//	cryptoIDInterface, ok := cryptoData["id"].(float64)
//	if !ok || cryptoIDInterface == 0 {
//		return 0, fmt.Errorf("cryptocurrency ID not found for symbol: %s", symbol)
//	}
//
//	cryptoID := int(cryptoIDInterface)
//	currentPrice, ok := priceData["price"].(float64)
//	if !ok {
//		return 0, fmt.Errorf("failed to convert price to float64")
//	}
//
//	if currentPrice >= targetPrice {
//		return currentPrice, fmt.Errorf("alert: %s has reached your target price of $%.2f. Current price: $%.2f", symbol, targetPrice, currentPrice)
//	}
//
//	notification := &models.PriceNotification{
//		CryptoID:    cryptoID,
//		Crypto:      symbol,
//		TargetPrice: targetPrice,
//		Username:    user.Username,
//		AskedAt:     time.Now().Format(time.RFC3339),
//		Status:      "Pending",
//	}
//
//	notificationRepo := storage.NewPGNotificationRepository(repo.conn)
//	if err := notificationRepo.SavePriceNotification(repo.conn, notification); err != nil {
//		return 0, fmt.Errorf("error saving notification: %v", err)
//	}
//
//	return currentPrice, fmt.Errorf("%s is still below your target price. Current price: $%.2f. Notification created.", symbol, currentPrice)
//}
//
//func CheckCryptocurrencyExists(symbol string) (bool, error) {
//	params := map[string]string{
//		"symbol":  symbol,
//		"convert": "USD",
//	}
//
//	response := api.GetAPIResponse("/quotes/latest", params)
//	var result map[string]interface{}
//	if err := json.Unmarshal(response, &result); err != nil {
//		return false, fmt.Errorf("error unmarshalling API response: %v", err)
//	}
//
//	data, dataOk := result["data"].(map[string]interface{})
//	if !dataOk || data[symbol] == nil {
//		return false, nil
//	}
//
//	return true, nil
//}

package repositories

import (
	"context"
	"cryptotracker/internal/api"
	"cryptotracker/models"
	"cryptotracker/pkg/config"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"math"
	"strconv"
	"strings"
	"time"
)

type CryptoRepository interface {
	DisplayTopCryptocurrencies(count int) ([]interface{}, error)
	SearchCryptocurrency(user *models.User, cryptoSymbol string) (float64, string, string, *models.Cryptocurrency, error)
	SetPriceAlert(user *models.User, symbol string, targetPrice float64) (float64, error)
}

type PostgresCryptoRepository struct {
	conn *pgx.Conn
}

func NewPostgresCryptoRepository(conn *pgx.Conn) *PostgresCryptoRepository {
	return &PostgresCryptoRepository{
		conn: conn,
	}
}

func (repo *PostgresCryptoRepository) DisplayTopCryptocurrencies(count int) ([]interface{}, error) {
	params := map[string]string{
		"start":   "1",
		"limit":   strconv.Itoa(count),
		"convert": "USD",
	}

	apiResponse := api.CoinMarketCapClient{}

	response := apiResponse.GetAPIResponse("/listings/latest", params)

	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("error unmarshalling API response: %v", err)
	}

	data, ok := result["data"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("data not found in the response")
	}

	return data, nil
}

func (repo *PostgresCryptoRepository) SearchCryptocurrency(user *models.User, cryptoSymbol string) (float64, string, string, *models.Cryptocurrency, error) {
	cryptoSymbol = strings.ToLower(cryptoSymbol)

	crypt := &models.Cryptocurrency{}

	params := map[string]string{
		"start":   "1",
		"limit":   "5000",
		"convert": "USD",
	}

	apiResponse := api.CoinMarketCapClient{}

	response := apiResponse.GetAPIResponse("/listings/latest", params)

	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		return 0, "", "", crypt, fmt.Errorf("error unmarshalling API response: %v", err)
	}

	data, ok := result["data"].([]interface{})
	if !ok || len(data) == 0 {
		return 0, "", "", crypt, fmt.Errorf("data not found or empty in the response")
	}

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

		if strings.ToLower(symbol) == cryptoSymbol || strings.ToLower(name) == cryptoSymbol {
			quote, ok := crypto["quote"].(map[string]interface{})
			if !ok {
				return 0, "", "", crypt, fmt.Errorf("quote data not found for cryptocurrency")
			}

			usd, ok := quote["USD"].(map[string]interface{})
			if !ok {
				return 0, "", "", crypt, fmt.Errorf("USD data not found in the quote")
			}

			price, ok := usd["price"].(float64)
			if !ok {
				return 0, "", "", crypt, fmt.Errorf("could not retrieve or convert the price")
			}

			price = math.Round(price*100) / 100.0

			cryptocurrency := &models.Cryptocurrency{
				CMCRank:     int(crypto["cmc_rank"].(float64)), // Assuming cmc_rank is float64
				DateAdded:   crypto["date_added"].(string),
				ID:          int(crypto["id"].(float64)), // Assuming id is float64
				LastUpdated: crypto["last_updated"].(string),
				Name:        name,
				Slug:        crypto["slug"].(string),
				Symbol:      symbol,
				Quote: models.Quote{
					USD: models.QuoteUSD{
						FullyDilutedMarketCap: usd["fully_diluted_market_cap"].(float64),
						PercentChange1H:       usd["percent_change_1h"].(float64),
						PercentChange24H:      usd["percent_change_24h"].(float64),
						PercentChange30D:      usd["percent_change_30d"].(float64),
						PercentChange60D:      usd["percent_change_60d"].(float64),
						PercentChange7D:       usd["percent_change_7d"].(float64),
						PercentChange90D:      usd["percent_change_90d"].(float64),
						Price:                 price,
					},
				},
			}

			return price, name, symbol, cryptocurrency, nil
		}
	}

	request := &models.UnavailableCryptoRequest{
		CryptoSymbol:   cryptoSymbol,
		UserName:       user.Username,
		RequestMessage: "Please add this cryptocurrency.",
		Status:         "Pending",
		Timestamp:      time.Now(),
	}

	if repo.conn == nil {
		return 0, "", "", crypt, fmt.Errorf("database connection is nil")
	}

	//unavailableCryptoRepo := storage.NewPGUnavailableCryptoRequestRepository(repo.conn)

	if err := repo.SaveUnavailableCryptoRequest(repo.conn, request); err != nil {
		return 0, "", "", crypt, fmt.Errorf("error saving unavailable crypto request: %v", err)
	}

	return 0, "", "", crypt, fmt.Errorf("request to add the cryptocurrency has been submitted")
}

func CheckCryptocurrencyExists(symbol string) (bool, float64, error) {
	params := map[string]string{
		"start":   "1",
		"limit":   "5000",
		"convert": "USD",
	}
	apiResponse := api.CoinMarketCapClient{}
	response := apiResponse.GetAPIResponse("/listings/latest", params)

	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		return false, 0, fmt.Errorf("error unmarshalling API response: %v", err)
	}

	data, ok := result["data"].([]interface{})
	if !ok || len(data) == 0 {
		return false, 0, fmt.Errorf("data not found or empty in the response")
	}

	symbol = strings.ToLower(symbol)

	for _, item := range data {
		crypto, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		cryptoSymbol, ok := crypto["symbol"].(string)
		if !ok {
			continue
		}

		if strings.ToLower(cryptoSymbol) == symbol {
			quote, ok := crypto["quote"].(map[string]interface{})
			if !ok {
				return true, 0, fmt.Errorf("quote data not found for cryptocurrency")
			}
			usd, ok := quote["USD"].(map[string]interface{})
			if !ok {
				return true, 0, fmt.Errorf("USD data not found in the quote")
			}
			price, ok := usd["price"].(float64)
			if !ok {
				return true, 0, fmt.Errorf("could not retrieve or convert the price")
			}
			return true, price, nil
		}
	}

	return false, 0, nil
}

func (repo *PostgresCryptoRepository) SetPriceAlert(user *models.User, symbol string, targetPrice float64) (float64, error) {
	exists, currentPrice, err := CheckCryptocurrencyExists(symbol)
	if err != nil {
		return 0, fmt.Errorf("error checking cryptocurrency: %v", err)
	}

	if !exists {
		return 0, fmt.Errorf("cryptocurrency with symbol %s does not exist", symbol)
	}

	if currentPrice >= targetPrice {
		return currentPrice, fmt.Errorf("alert: %s has reached your target price of $%.2f. Current price: $%.2f", symbol, targetPrice, currentPrice)
	}

	notification := &models.PriceNotification{
		CryptoID:    0, // We don't have the ID anymore, consider removing this field if not needed
		Crypto:      symbol,
		TargetPrice: targetPrice,
		Username:    user.Username,
		AskedAt:     time.Now().Format(time.RFC3339),
		Status:      "Pending",
	}

	notificationRepository := NewPostgresNotificationRepository(repo.conn)
	if err := notificationRepository.SavePriceNotification(repo.conn, notification); err != nil {
		return 0, fmt.Errorf("error saving notification: %v", err)
	}

	return currentPrice, fmt.Errorf("%s is still below your target price. Current price: $%.2f. Notification created.", symbol, currentPrice)
}

func (repo *PostgresCryptoRepository) SaveUnavailableCryptoRequest(conn *pgx.Conn, request *models.UnavailableCryptoRequest) error {
	columns := []string{"crypto_symbol", "username", "request_message", "status", "timestamp"}
	query, err := config.BuildInsertQuery("unavailable_cryptos", columns)
	if err != nil {
		return fmt.Errorf("failed to build insert query: %v", err)
	}

	request.Timestamp = time.Now()

	_, err = conn.Exec(context.Background(), query, request.CryptoSymbol, request.UserName, request.RequestMessage, request.Status, request.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to save unavailable crypto request: %v", err)
	}

	return nil
}
