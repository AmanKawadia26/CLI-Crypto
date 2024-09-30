package handlers

import (
	"bytes"
	errs "cryptotracker/REST-API/errors"
	Handlers "cryptotracker/REST-API/handlers"
	mock_middleware "cryptotracker/test/REST-API/mocks"
	mock_services "cryptotracker/test/internal/mocks/services"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Setup the necessary mocks and handlers
func setupAuthTest(t *testing.T) (*Handlers.AuthHandler, *mock_services.MockAuthService, *mock_middleware.MockTokenService, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	authService := mock_services.NewMockAuthService(ctrl)
	tokenService := mock_middleware.NewMockTokenService(ctrl)
	authHandler := Handlers.NewAuthHandler(authService)
	return authHandler, authService, tokenService, ctrl
}

//func TestLoginHandler_Success(t *testing.T) {
//	authHandler, authService, tokenService, ctrl := setupAuthTest(t)
//	defer ctrl.Finish()
//
//	loginReq := Handlers.LoginRequest{
//		Username: "testuser",
//		Password: "password123",
//	}
//
//	// Mock the Login service response
//	user := &models.User{
//		Username: "testuser",
//		Role:     "user",
//	}
//	authService.EXPECT().Login(loginReq.Username, loginReq.Password).Return(user, "token123", nil)
//
//	// Mock the Token generation
//	tokenService.EXPECT().GenerateToken(loginReq.Username, user.Role).Return("valid-token", nil)
//
//	reqBody, _ := json.Marshal(loginReq)
//	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
//	rec := httptest.NewRecorder()
//
//	// Call the handler function
//	authHandler.LoginHandler(rec, req)
//
//	// Check if status code is OK
//	assert.Equal(t, http.StatusOK, rec.Code)
//
//	// Check the response body
//	expectedResp := response.Response{
//		Status:  "success",
//		Message: "Login successful",
//	}
//	var actualResp response.Response
//	json.Unmarshal(rec.Body.Bytes(), &actualResp)
//	assert.Equal(t, expectedResp.Status, actualResp.Status)
//	assert.Equal(t, expectedResp.Message, actualResp.Message)
//}

func TestLoginHandler_Failure_InvalidCredentials(t *testing.T) {
	authHandler, authService, _, ctrl := setupAuthTest(t)
	defer ctrl.Finish()

	loginReq := Handlers.LoginRequest{
		Username: "wronguser",
		Password: "wrongpassword",
	}

	// Mock invalid login credentials
	authService.EXPECT().Login(loginReq.Username, loginReq.Password).Return(nil, "", errs.NewInvalidCredentialsError())

	reqBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	rec := httptest.NewRecorder()

	// Call the handler function
	authHandler.LoginHandler(rec, req)

	// Check if status code is Unauthorized
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

//func TestSignupHandler_Success(t *testing.T) {
//	authHandler, authService, _, ctrl := setupAuthTest(t)
//	defer ctrl.Finish()
//
//	signupReq := Handlers.SignupRequest{
//		Username: "newuser",
//		Password: "password123",
//		Email:    "newuser@example.com",
//		Mobile:   1234567890,
//	}
//
//	// Mock the Signup service response
//	authService.EXPECT().Signup(gomock.Any()).Return(nil)
//
//	reqBody, _ := json.Marshal(signupReq)
//	req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(reqBody))
//	rec := httptest.NewRecorder()
//
//	// Call the handler function
//	authHandler.SignupHandler(rec, req)
//
//	// Check if status code is OK
//	assert.Equal(t, http.StatusOK, rec.Code)
//
//	// Check the response body
//	expectedResp := response.Response{
//		Status:  "success",
//		Message: "Signup successful",
//	}
//	var actualResp response.Response
//	json.Unmarshal(rec.Body.Bytes(), &actualResp)
//	assert.Equal(t, expectedResp.Status, actualResp.Status)
//	assert.Equal(t, expectedResp.Message, actualResp.Message)
//}
//
//func TestLogoutHandler_Success(t *testing.T) {
//	authHandler, _, tokenService, ctrl := setupAuthTest(t)
//	defer ctrl.Finish()
//
//	// Mock the token extraction and parsing
//	tokenService.EXPECT().ExtractToken(gomock.Any()).Return("valid-token")
//	tokenService.EXPECT().ParseToken("valid-token").Return(&middleware.Claims{
//		Username: "testuser",
//	}, nil)
//	tokenService.EXPECT().BlacklistToken("valid-token")
//
//	req := httptest.NewRequest("POST", "/logout", nil)
//	rec := httptest.NewRecorder()
//
//	// Call the handler function
//	authHandler.LogoutHandler(rec, req)
//
//	// Check if status code is OK
//	assert.Equal(t, http.StatusOK, rec.Code)
//}
