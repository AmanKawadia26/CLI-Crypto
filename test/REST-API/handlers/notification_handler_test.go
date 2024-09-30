package handlers

import (
	"context"
	errs "cryptotracker/REST-API/errors"
	Handlers "cryptotracker/REST-API/handlers"
	"cryptotracker/internal/services"
	mock_services "cryptotracker/test/internal/mocks/services"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckNotificationHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock NotificationService
	mockNotificationService := mock_services.NewMockNotificationService(ctrl)

	// Set up the mock response
	mockNotifications := []services.Notification{
		{Message: "Price alert met"},
		{Message: "Crypto request approved"},
	}

	// Expect the CheckNotification function to be called with the correct username
	mockNotificationService.EXPECT().CheckNotification("testuser").Return(mockNotifications, nil)

	// Create the handler with the mock service
	handler := Handlers.NewNotificationHandler(mockNotificationService)

	// Create a request with a valid username context
	req, err := http.NewRequest("GET", "/notifications", nil)
	assert.NoError(t, err)

	// Set username in the context
	ctx := context.WithValue(req.Context(), "username", "testuser")
	req = req.WithContext(ctx)

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler.CheckNotificationHandler(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response body
	var jsonResponse map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &jsonResponse)
	assert.NoError(t, err)

	// Check the response content
	assert.Equal(t, "success", jsonResponse["status"])
	assert.Equal(t, "Notifications retrieved", jsonResponse["message"])
	assert.Len(t, jsonResponse["data"], 2)
}

func TestCheckNotificationHandler_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock NotificationService (unused in this test)
	mockNotificationService := mock_services.NewMockNotificationService(ctrl)

	// Create the handler with the mock service
	handler := Handlers.NewNotificationHandler(mockNotificationService)

	// Create a request without a valid username context
	req, err := http.NewRequest("GET", "/notifications", nil)
	assert.NoError(t, err)

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler.CheckNotificationHandler(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	// Parse the response body
	var jsonResponse errs.AppError
	err = json.Unmarshal(rr.Body.Bytes(), &jsonResponse)
	assert.NoError(t, err)

	// Check the error response content
	assert.Equal(t, errs.NewUnauthorizedAccessError().Message, jsonResponse.Message)
}

func TestCheckNotificationHandler_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock NotificationService
	mockNotificationService := mock_services.NewMockNotificationService(ctrl)

	// Expect the CheckNotification function to return an error
	mockNotificationService.EXPECT().CheckNotification("testuser").Return(nil, errs.NewFailedToCheckNotificationsError())

	// Create the handler with the mock service
	handler := Handlers.NewNotificationHandler(mockNotificationService)

	// Create a request with a valid username context
	req, err := http.NewRequest("GET", "/notifications", nil)
	assert.NoError(t, err)

	// Set username in the context
	ctx := context.WithValue(req.Context(), "username", "testuser")
	req = req.WithContext(ctx)

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler.CheckNotificationHandler(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Parse the response body
	var jsonResponse errs.AppError
	err = json.Unmarshal(rr.Body.Bytes(), &jsonResponse)
	assert.NoError(t, err)

	// Check the error response content
	assert.Equal(t, errs.NewFailedToCheckNotificationsError().Message, jsonResponse.Message)
}
