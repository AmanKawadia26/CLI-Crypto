package services

import (
	//"cryptotracker/internal/notification"
	"cryptotracker/internal/repositories"
	"cryptotracker/models"
	"cryptotracker/pkg/utils"
	"github.com/pkg/errors"
	//"github.com/jackc/pgx/v4"
)

type AuthService interface {
	Login(username, password string) (*models.User, string, error)
	Signup(user *models.User) error
}

type AuthServiceImpl struct {
	repo                repositories.AuthRepository
	NotificationService NotificationService
}

func NewAuthService(repo repositories.AuthRepository, notificationService NotificationService) AuthService {
	return &AuthServiceImpl{repo: repo,
		NotificationService: notificationService}
}

func (s *AuthServiceImpl) Login(username, password string) (*models.User, string, error) {

	// Call the repository function to fetch the user from the database
	user, err := s.repo.LoginDBRepository(username)
	if err != nil {
		return nil, "", err // Return the error if user is not found or any other DB issue
	}

	// Compare the provided password with the stored hashed password
	hashedPassword := utils.HashPassword(password)
	if user.Password != hashedPassword {
		return nil, "", errors.New("invalid username or password")
	}

	// Check and display any notifications for the user
	//s.NotificationService.CheckNotification(username)
	//notification.CheckNotification(s.repo.conn, username)

	// Return the user object and role
	return user, user.Role, nil
}

func (s *AuthServiceImpl) Signup(user *models.User) error {
	err := s.repo.SignupDBRepository(user)
	return err
}
