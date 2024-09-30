package handlers

import (
	Handlers "cryptotracker/REST-API/handlers"
	"cryptotracker/models"
	mock_services "cryptotracker/test/internal/mocks/services"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

//func TestNewAdminHandler(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockAdminService := mock_services.NewMockAdminService(ctrl)
//
//	handler := Handlers.NewAdminHandler(mockAdminService)
//	assert.NotNil(t, handler)
//	assert.Equal(t, mockAdminService, handler.)
//}

// TestProfiles tests the AdminHandler's Profiles method
func TestProfiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdminService := mock_services.NewMockAdminService(ctrl)
	adminHandler := Handlers.NewAdminHandler(mockAdminService)

	tests := []struct {
		name         string
		username     string
		serviceSetup func()
		expectedCode int
		expectedBody string
	}{
		{
			name:     "Success - Fetch all profiles",
			username: "",
			serviceSetup: func() {
				mockAdminService.EXPECT().ViewUserProfiles().Return([]*models.User{
					{Username: "user1"}, {Username: "user2"},
				}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `"status":"success","message":"User profiles fetched successfully"`,
		},
		{
			name:     "Success - Fetch specific profile",
			username: "user1",
			serviceSetup: func() {
				mockAdminService.EXPECT().ViewUserProfiles().Return([]*models.User{
					{Username: "user1"},
				}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `"status":"success","message":"User profile fetched successfully"`,
		},
		{
			name:     "Error - User not found",
			username: "nonexistent",
			serviceSetup: func() {
				mockAdminService.EXPECT().ViewUserProfiles().Return([]*models.User{
					{Username: "user1"},
				}, nil)
			},
			expectedCode: http.StatusNotFound,
			expectedBody: `"message":"User not found"`,
		},
		{
			name:     "Error - Service failure",
			username: "",
			serviceSetup: func() {
				mockAdminService.EXPECT().ViewUserProfiles().Return(nil, errors.New("service error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `"message":"Failed to fetch user profiles"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.serviceSetup()

			req := httptest.NewRequest(http.MethodGet, "/admin/profiles", nil)
			if tt.username != "" {
				q := req.URL.Query()
				q.Add("username", tt.username)
				req.URL.RawQuery = q.Encode()
			}

			rec := httptest.NewRecorder()
			adminHandler.Profiles(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tt.expectedBody)
		})
	}
}

// TestDeleteUser tests the AdminHandler's DeleteUser method
func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdminService := mock_services.NewMockAdminService(ctrl)
	adminHandler := Handlers.NewAdminHandler(mockAdminService)

	tests := []struct {
		name         string
		username     string
		serviceSetup func()
		expectedCode int
		expectedBody string
	}{
		{
			name:     "Success - Delete user",
			username: "user1",
			serviceSetup: func() {
				mockAdminService.EXPECT().DeleteUser("user1").Return(nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `"status":"success","message":"User deleted successfully"`,
		},
		{
			name:     "Error - Service failure",
			username: "user1",
			serviceSetup: func() {
				mockAdminService.EXPECT().DeleteUser("user1").Return(errors.New("service error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `"message":"Failed to delete user"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.serviceSetup()

			req := httptest.NewRequest(http.MethodDelete, "/admin/user/"+tt.username, nil)
			req = mux.SetURLVars(req, map[string]string{"username": tt.username})
			rec := httptest.NewRecorder()
			adminHandler.DeleteUser(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tt.expectedBody)
		})
	}
}

// TestDelegateUser tests the AdminHandler's DelegateUser method
func TestDelegateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdminService := mock_services.NewMockAdminService(ctrl)
	adminHandler := Handlers.NewAdminHandler(mockAdminService)

	tests := []struct {
		name         string
		username     string
		serviceSetup func()
		expectedCode int
		expectedBody string
	}{
		{
			name:     "Success - Delegate user",
			username: "user1",
			serviceSetup: func() {
				mockAdminService.EXPECT().ChangeUserStatus("user1").Return(nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `"status":"success","message":"User delegated to admin successfully"`,
		},
		{
			name:     "Error - Service failure",
			username: "user1",
			serviceSetup: func() {
				mockAdminService.EXPECT().ChangeUserStatus("user1").Return(errors.New("service error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `"message":"Failed to delegate user"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.serviceSetup()

			req := httptest.NewRequest(http.MethodPost, "/admin/user/delegate/"+tt.username, nil)
			req = mux.SetURLVars(req, map[string]string{"username": tt.username})
			rec := httptest.NewRecorder()
			adminHandler.DelegateUser(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tt.expectedBody)
		})
	}
}

// Helper function to prepare test cases for unavailable crypto requests
func prepareUnavailableCryptoRequests() []*models.UnavailableCryptoRequest {
	return []*models.UnavailableCryptoRequest{
		{CryptoSymbol: "BTC", UserName: "user1"},
		{CryptoSymbol: "ETH", UserName: "user2"},
	}
}

// TestUnavailableCryptoRequests tests the UnavailableCryptoRequests handler
func TestUnavailableCryptoRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdminService := mock_services.NewMockAdminService(ctrl)
	adminHandler := Handlers.NewAdminHandler(mockAdminService)

	tests := []struct {
		name         string
		cryptoSymbol string
		serviceSetup func()
		expectedCode int
		expectedBody string
	}{
		{
			name: "Success - Fetch all crypto requests",
			serviceSetup: func() {
				mockAdminService.EXPECT().ManageUserRequests().Return(prepareUnavailableCryptoRequests(), nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `"status":"success","message":"Requests fetched successfully"`,
		},
		{
			name: "Error - Service failure",
			serviceSetup: func() {
				mockAdminService.EXPECT().ManageUserRequests().Return(nil, errors.New("service error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `"message":"Failed to fetch requests"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.serviceSetup()

			req := httptest.NewRequest(http.MethodGet, "/admin/requests", nil)
			rec := httptest.NewRecorder()
			adminHandler.UnavailableCryptoRequests(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tt.expectedBody)
		})
	}
}
