package services

import (
	"cryptotracker/internal/repositories"
	"cryptotracker/models"
)

type UserServices interface {
	GetUserProfile(username string) (*models.User, error)
}

type UserServiceImpl struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{repo: repo}
}

// Business logic to get the user profile
func (s *UserServiceImpl) GetUserProfile(username string) (*models.User, error) {
	return s.repo.GetUserProfile(username)
}
