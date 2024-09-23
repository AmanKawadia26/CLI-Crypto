package api_test

import (
	//"bytes"
	"cryptotracker/internal/api"
	"cryptotracker/pkg/config"
	//"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAPIResponse(t *testing.T) {
	// Set up a test server with a mock response
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request method and URL
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/cryptocurrency", r.URL.Path)
		assert.Equal(t, "test_key", r.Header.Get("X-CMC_PRO_API_KEY"))

		// Set up mock response
		response := []byte(`{"status": "success", "data": {}}`)
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(response)
		require.NoError(t, err)
	}))
	defer mockServer.Close()

	// Override the baseURL with the mock server URL
	originalBaseURL := api.GetBaseURL()
	api.SetBaseURL(mockServer.URL)

	// Set up the config for the test
	config.AppConfig.APIKey = "test_key"

	// Define parameters for the test
	params := map[string]string{"param1": "value1", "param2": "value2"}

	// Call the function
	response := api.GetAPIResponse("/v1/cryptocurrency", params)

	// Validate the response
	expectedResponse := `{"status": "success", "data": {}}`
	assert.Equal(t, expectedResponse, string(response))

	// Restore the original baseURL
	api.SetBaseURL(originalBaseURL)
}
