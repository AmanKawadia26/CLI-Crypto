package response

import (
	"encoding/json"
	"net/http"
)

// General response struct for all API responses
type Response struct {
	//StatusCode int         `json:"status_code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`  // Optional data field
	Token   string      `json:"token,omitempty"` // Optional token field for login responses
}

// SendJSONResponse writes the response to the client with the appropriate status code
func SendJSONResponse(w http.ResponseWriter, statusCode int, status string, message string, data interface{}, token string) {
	response := Response{
		//StatusCode: statusCode,
		Status:  status,
		Message: message,
		Data:    data,
		Token:   token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
