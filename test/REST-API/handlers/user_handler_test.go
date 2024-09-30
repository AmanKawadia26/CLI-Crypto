package handlers

import (
	"context"
	Handlers "cryptotracker/REST-API/handlers"
	"cryptotracker/models"
	mock_services "cryptotracker/test/internal/mocks/services"
	"errors"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestUserProfile_Success tests the case where the user profile is successfully retrieved
func TestUserProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_services.NewMockUserServices(ctrl)

	// Mock User Profile
	mockUserProfile := &models.User{
		Email:    "testuser@example.com",
		Mobile:   1234567890,
		Username: "testuser",
	}

	// Expect GetUserProfile to be called and return the mockUserProfile without error
	mockUserService.EXPECT().GetUserProfile("testuser").Return(mockUserProfile, nil)

	// Create a UserHandler instance
	userHandler := Handlers.NewUserHandler(mockUserService)

	// Create a request with a context that includes the username
	req, err := http.NewRequest("GET", "/users/me", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v", err)
	}
	ctx := context.WithValue(req.Context(), "username", "testuser")
	req = req.WithContext(ctx)

	// Record the response
	rr := httptest.NewRecorder()

	// Call the handler
	userHandler.UserProfile(rr, req)

	// Check the response code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := `{"status":"success","message":"User profile retrieved successfully","data":{"email":"testuser@example.com","mobile":"1234567890","username":"testuser"}}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// TestUserProfile_Unauthorized tests the case where the username is missing from the request context
func TestUserProfile_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_services.NewMockUserServices(ctrl)

	// Create a UserHandler instance
	userHandler := Handlers.NewUserHandler(mockUserService)

	// Create a request without a username in the context
	req, err := http.NewRequest("GET", "/users/me", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v", err)
	}

	// Record the response
	rr := httptest.NewRecorder()

	// Call the handler
	userHandler.UserProfile(rr, req)

	// Check the response code
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	// Check the response body for unauthorized error
	expected := `{"code":4012,"message":"Unauthorized access"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// TestUserProfile_FailedToRetrieve tests the case where the service fails to retrieve the user profile
func TestUserProfile_FailedToRetrieve(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_services.NewMockUserServices(ctrl)

	// Expect GetUserProfile to be called and return an error
	mockUserService.EXPECT().GetUserProfile("testuser").Return(nil, errors.New("some error"))

	// Create a UserHandler instance
	userHandler := Handlers.NewUserHandler(mockUserService)

	// Create a request with a context that includes the username
	req, err := http.NewRequest("GET", "/users/me", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v", err)
	}
	ctx := context.WithValue(req.Context(), "username", "testuser")
	req = req.WithContext(ctx)

	// Record the response
	rr := httptest.NewRecorder()

	// Call the handler
	userHandler.UserProfile(rr, req)

	// Check the response code
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	// Check the response body
	expected := `{"code":5012,"message":"Failed to retrieve user profile"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
