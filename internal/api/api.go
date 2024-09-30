//package api
//
//import (
//	"cryptotracker/pkg/config"
//	"io/ioutil"
//	"log"
//	"net/http"
//	"time"
//)
//
//var baseURL = "https://pro-api.coinmarketcap.com/v1/cryptocurrency"
//
//func GetBaseURL() string {
//	return baseURL
//}
//
//func SetBaseURL(newBaseURL string) {
//	baseURL = newBaseURL
//}
//
//func GetAPIResponse(endpoint string, params map[string]string) []byte {
//	client := &http.Client{Timeout: 30 * time.Second}
//	req, _ := http.NewRequest("GET", GetBaseURL()+endpoint, nil)
//
//	req.Header.Add("X-CMC_PRO_API_KEY", config.AppConfig.APIKey)
//
//	q := req.URL.Query()
//	for key, value := range params {
//		q.Add(key, value)
//	}
//	req.URL.RawQuery = q.Encode()
//
//	resp, err := client.Do(req)
//	if err != nil {
//		log.Fatalf("Error making API request: %v", err)
//	}
//	defer resp.Body.Close()
//
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		log.Fatalf("Error reading API response: %v", err)
//	}
//
//	return body
//}

package api

import (
	"cryptotracker/pkg/config"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var baseURL = "https://pro-api.coinmarketcap.com/v1/cryptocurrency"

// APIClient defines the interface for making API requests.
// This interface allows for easier testing and flexibility.
type APIClient interface {
	GetAPIResponse(endpoint string, params map[string]string) []byte
}

// CoinMarketCapClient implements the APIClient interface.
// It provides methods to interact with the CoinMarketCap API.
type CoinMarketCapClient struct{}

// GetBaseURL returns the base URL for the CoinMarketCap API.
func GetBaseURL() string {
	return baseURL
}

// SetBaseURL allows updating the base URL for the CoinMarketCap API.
func SetBaseURL(newBaseURL string) {
	baseURL = newBaseURL
}

// GetAPIResponse makes an HTTP GET request to the given API endpoint with provided parameters.
// This method implements the APIClient interface.
func (c *CoinMarketCapClient) GetAPIResponse(endpoint string, params map[string]string) []byte {
	client := &http.Client{Timeout: 30 * time.Second}
	req, _ := http.NewRequest("GET", GetBaseURL()+endpoint, nil)

	// Set the API key header
	req.Header.Add("X-CMC_PRO_API_KEY", config.AppConfig.APIKey)

	// Add query parameters to the URL
	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	// Make the request and handle any errors
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making API request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading API response: %v", err)
	}

	return body
}
