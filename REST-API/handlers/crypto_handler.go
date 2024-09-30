package Handlers

import (
	"cryptotracker/REST-API/errors" // Import the errors package
	"cryptotracker/REST-API/response"
	"cryptotracker/internal/services"
	"cryptotracker/models"
	"cryptotracker/pkg/logger"
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type CryptoHandler struct {
	cryptoService services.CryptoService
}

func NewCryptoHandler(cryptoService services.CryptoService) *CryptoHandler {
	return &CryptoHandler{
		cryptoService: cryptoService,
	}
}

func (h *CryptoHandler) DisplayTopCryptos(w http.ResponseWriter, r *http.Request) {
	count := r.URL.Query().Get("count")
	if count == "" {
		count = "10"
	}

	countInt, err := strconv.Atoi(count)
	if err != nil {
		logger.Logger.Error("Invalid count parameter", err)
		errs.NewInvalidRequestPayloadError().ToJSON(w)
		return
	}

	cryptos, err := h.cryptoService.DisplayTopCryptocurrencies(countInt)
	if err != nil {
		logger.Logger.Error("Error retrieving top cryptocurrencies", err)
		errs.NewRetrievingCryptosError().ToJSON(w)
		return
	}

	// Create a slice of Cryptocurrency_Response instead of a map
	filteredCryptos := make([]models.Cryptocurrency, 0)
	for _, crypto := range cryptos {
		cryptoMap, ok := crypto.(map[string]interface{})
		if !ok {
			logger.Logger.Error("Failed to type assert crypto data")
			errs.NewRetrievingCryptosError().ToJSON(w)
			return
		}

		// Extract relevant fields using the map
		quote, ok := cryptoMap["quote"].(map[string]interface{})
		if !ok {
			logger.Logger.Error("Failed to extract quote from crypto data")
			errs.NewRetrievingCryptosError().ToJSON(w)
			return
		}

		usdQuote, ok := quote["USD"].(map[string]interface{})
		if !ok {
			logger.Logger.Error("Failed to extract USD quote from crypto data")
			errs.NewRetrievingCryptosError().ToJSON(w)
			return
		}
		filteredCrypto := models.Cryptocurrency{
			CMCRank:     int(cryptoMap["cmc_rank"].(float64)), // Cast float64 to int
			DateAdded:   cryptoMap["date_added"].(string),
			ID:          int(cryptoMap["id"].(float64)), // Cast float64 to int
			LastUpdated: cryptoMap["last_updated"].(string),
			Name:        cryptoMap["name"].(string),
			Slug:        cryptoMap["slug"].(string),
			Symbol:      cryptoMap["symbol"].(string),
			Quote: models.Quote{
				USD: models.QuoteUSD{
					FullyDilutedMarketCap: usdQuote["fully_diluted_market_cap"].(float64),
					PercentChange1H:       usdQuote["percent_change_1h"].(float64),
					PercentChange24H:      usdQuote["percent_change_24h"].(float64),
					PercentChange30D:      usdQuote["percent_change_30d"].(float64),
					PercentChange60D:      usdQuote["percent_change_60d"].(float64),
					PercentChange7D:       usdQuote["percent_change_7d"].(float64),
					PercentChange90D:      usdQuote["percent_change_90d"].(float64),
					Price:                 math.Round(usdQuote["price"].(float64)*100) / 100.0,
				},
			},
		}
		filteredCryptos = append(filteredCryptos, filteredCrypto)
	}

	logger.Logger.Info("Top cryptocurrencies retrieved successfully")
	response.SendJSONResponse(w, http.StatusOK, "success", "Top cryptocurrencies retrieved successfully", filteredCryptos, "")
}

func (h *CryptoHandler) DisplayCryptoByName(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value("username").(string)
	if !ok || username == "" {
		logger.Logger.Error("Unauthorized access")
		errs.NewUnauthorizedAccessError().ToJSON(w)
		return
	}

	cryptoSymbol := strings.TrimPrefix(r.URL.Path, "/cryptos/")
	if cryptoSymbol == "" {
		logger.Logger.Warn("No cryptocurrency symbol provided")
		errs.NewMissingCryptoSymbolError().ToJSON(w)
		return
	}

	user := &models.User{Username: username}

	_, _, _, crypt, err := h.cryptoService.SearchCryptocurrency(user, cryptoSymbol)

	if err != nil && strings.Contains(err.Error(), "to add the cryptocurrency has been submitted") {
		logger.Logger.Info("Request to add the cryptocurrency has been submitted", err)
		response.SendJSONResponse(w, http.StatusOK, "success", "Request to add "+cryptoSymbol+" cryptocurrency has been added", nil, "")
		return
	} else if err != nil {
		logger.Logger.Error("Error searching for cryptocurrency", err)
		errs.NewSearchingCryptoError().ToJSON(w)
		return
	}

	// Only return the `crypt` (Cryptocurrency struct) in the response
	logger.Logger.Info("Cryptocurrency details retrieved successfully")
	response.SendJSONResponse(w, http.StatusOK, "success", "Cryptocurrency details retrieved successfully", crypt, "")
}

var AlertReq struct {
	Symbol      string  `json:"crypto_symbol"`
	TargetPrice float64 `json:"target_price"`
}

func (h *CryptoHandler) SetPriceAlert(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value("username").(string)
	if !ok || username == "" {
		logger.Logger.Error("Unauthorized access")
		errs.NewUnauthorizedAccessError().ToJSON(w)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&AlertReq)
	if err != nil {
		logger.Logger.Error("Failed to decode price alert request", err)
		errs.NewInvalidRequestPayloadError().ToJSON(w)
		return
	}

	user := &models.User{Username: username}

	price, err := h.cryptoService.SetPriceAlert(user, AlertReq.Symbol, AlertReq.TargetPrice)

	if price == 0 {
		logger.Logger.Error("Error setting price alert", err)
		errs.NewSettingPriceAlertError().ToJSON(w)
		return
	} else if price >= AlertReq.TargetPrice {
		logger.Logger.Warn("Current price is higher than requested", err)
		errs.NewCurrentPriceHigherError().ToJSON(w)
		return
	}

	logger.Logger.Info("Price Alert Request created successfully")
	response.SendJSONResponse(w, http.StatusOK, "success", "Price alert set successfully", map[string]interface{}{
		"current_price": price,
	}, "")
}
