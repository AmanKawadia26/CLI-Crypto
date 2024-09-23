package services_test

import (
	"cryptotracker/internal/services"
	"cryptotracker/models"
	mock_repositories "cryptotracker/test/internal/mocks/repository"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUserProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock UserRepository
	mockRepo := mock_repositories.NewMockUserRepository(ctrl)

	// Define the expected user profile
	expectedUser := &models.User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Mobile:   1234567890,
		//PAN:      "ABCDE1234F",
	}

	// Set up the expectation for GetUserProfile
	mockRepo.EXPECT().
		GetUserProfile("testuser").
		Return(expectedUser, nil).
		Times(1)

	// Create UserService with the mock repository
	service := services.NewUserService(mockRepo)

	// Call the method under test
	result, err := service.GetUserProfile("testuser")

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, result)
}

func TestGetUserProfile_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock UserRepository
	mockRepo := mock_repositories.NewMockUserRepository(ctrl)

	// Set up the expectation for GetUserProfile to return an error
	mockRepo.EXPECT().
		GetUserProfile("testuser").
		Return(nil, assert.AnError).
		Times(1)

	// Create UserService with the mock repository
	service := services.NewUserService(mockRepo)

	// Call the method under test
	result, err := service.GetUserProfile("testuser")

	// Assert the results
	assert.Error(t, err)
	assert.Nil(t, result)
}
