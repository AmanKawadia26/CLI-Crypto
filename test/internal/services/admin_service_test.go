package services_test

import (
	"cryptotracker/internal/services"
	"cryptotracker/models"
	mock_repositories "cryptotracker/test/internal/mocks/repository"
	"errors"
	"github.com/golang/mock/gomock"
	//"github.com/jackc/pgx/v4"
	"testing"
)

func TestAdminServiceImpl_ChangeUserStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repositories.NewMockAdminRepository(ctrl)
	service := services.NewAdminService(mockRepo)

	testCases := []struct {
		name     string
		username string
		mockErr  error
	}{
		{"Success", "testuser", nil},
		{"Error", "erroruser", errors.New("database error")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.EXPECT().ChangeUserStatus(gomock.Any(), tc.username).Return(tc.mockErr)
			err := service.ChangeUserStatus(nil, tc.username)
			if (err != nil) != (tc.mockErr != nil) {
				t.Errorf("Expected error: %v, got: %v", tc.mockErr, err)
			}
		})
	}
}

func TestAdminServiceImpl_DeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repositories.NewMockAdminRepository(ctrl)
	service := services.NewAdminService(mockRepo)

	testCases := []struct {
		name     string
		username string
		mockErr  error
	}{
		{"Success", "testuser", nil},
		{"Error", "erroruser", errors.New("database error")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.EXPECT().DeleteUser(gomock.Any(), tc.username).Return(tc.mockErr)
			err := service.DeleteUser(nil, tc.username)
			if (err != nil) != (tc.mockErr != nil) {
				t.Errorf("Expected error: %v, got: %v", tc.mockErr, err)
			}
		})
	}
}

func TestAdminServiceImpl_ViewUserProfiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repositories.NewMockAdminRepository(ctrl)
	service := services.NewAdminService(mockRepo)

	mockUsers := []*models.User{
		{Username: "user1", Email: "user1@example.com"},
		{Username: "user2", Email: "user2@example.com"},
	}

	testCases := []struct {
		name      string
		mockUsers []*models.User
		mockErr   error
	}{
		{"Success", mockUsers, nil},
		{"Error", nil, errors.New("database error")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.EXPECT().ViewUserProfiles(gomock.Any()).Return(tc.mockUsers, tc.mockErr)
			users, err := service.ViewUserProfiles(nil)
			if (err != nil) != (tc.mockErr != nil) {
				t.Errorf("Expected error: %v, got: %v", tc.mockErr, err)
			}
			if err == nil && len(users) != len(tc.mockUsers) {
				t.Errorf("Expected %d users, got %d", len(tc.mockUsers), len(users))
			}
		})
	}
}

func TestAdminServiceImpl_ManageUserRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repositories.NewMockAdminRepository(ctrl)
	service := services.NewAdminService(mockRepo)

	mockRequests := []*models.UnavailableCryptoRequest{
		{CryptoSymbol: "ABC", UserName: "user1"},
		{CryptoSymbol: "XYZ", UserName: "user2"},
	}

	testCases := []struct {
		name         string
		mockRequests []*models.UnavailableCryptoRequest
		mockErr      error
	}{
		{"Success", mockRequests, nil},
		{"Error", nil, errors.New("database error")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.EXPECT().ManageUserRequests(gomock.Any()).Return(tc.mockRequests, tc.mockErr)
			requests, err := service.ManageUserRequests(nil)
			if (err != nil) != (tc.mockErr != nil) {
				t.Errorf("Expected error: %v, got: %v", tc.mockErr, err)
			}
			if err == nil && len(requests) != len(tc.mockRequests) {
				t.Errorf("Expected %d requests, got %d", len(tc.mockRequests), len(requests))
			}
		})
	}
}

func TestAdminServiceImpl_UpdateRequestStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repositories.NewMockAdminRepository(ctrl)
	service := services.NewAdminService(mockRepo)

	request := &models.UnavailableCryptoRequest{
		CryptoSymbol: "ABC",
		UserName:     "user1",
		Status:       "pending",
	}

	testCases := []struct {
		name    string
		request *models.UnavailableCryptoRequest
		status  string
		mockErr error
	}{
		{"Success", request, "approved", nil},
		{"Error", request, "rejected", errors.New("database error")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.EXPECT().SaveUnavailableCryptoRequest(gomock.Any(), gomock.Any()).Return(tc.mockErr)
			err := service.UpdateRequestStatus(nil, tc.request, tc.status)
			if (err != nil) != (tc.mockErr != nil) {
				t.Errorf("Expected error: %v, got: %v", tc.mockErr, err)
			}
			if err == nil && tc.request.Status != tc.status {
				t.Errorf("Expected status to be '%s', got '%s'", tc.status, tc.request.Status)
			}
		})
	}
}

func TestNewAdminService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repositories.NewMockAdminRepository(ctrl)
	service := services.NewAdminService(mockRepo)

	if service == nil {
		t.Error("Expected non-nil AdminService")
	}

	_, ok := service.(*services.AdminServiceImpl)
	if !ok {
		t.Error("Expected AdminServiceImpl type")
	}
}
