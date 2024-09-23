package ui

import (
	"cryptotracker/internal/repositories"
	"cryptotracker/internal/services"
	"cryptotracker/models"
	"cryptotracker/pkg/utils"
	"cryptotracker/pkg/validation"
	//"database/sql"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/jackc/pgx/v4"
)

// AuthenticateUser handles the login/signup process
func (ui *UI) AuthenticateUser(conn *pgx.Conn) (*models.User, string) {
	for {
		ClearScreen()
		DisplayAuthMenu()

		authRepo := repositories.NewPostgresAuthRepository(conn)
		authService := services.NewAuthService(authRepo, ui.notificationService)

		var choice int
		color.New(color.FgCyan).Print("Enter your choice: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			if user, Role, err := ui.LoginUI(authService); err == nil {
				return user, Role
			} else {
				color.New(color.FgRed).Println("Login failed:", err)
			}
		case 2:
			if _, err := ui.SignupUI(conn, authService); err != nil {
				color.New(color.FgRed).Println("Signup failed:", err)
			} else {
				color.New(color.FgGreen).Println("Signup successful. Please login.")
			}
		case 3:
			color.New(color.FgYellow).Println("Exiting...")
			return nil, ""
		default:
			color.New(color.FgRed).Println("Invalid choice, please try again.")
		}
	}
}

// SignupUI handles user input and validation for PostgreSQL
func (ui *UI) SignupUI(conn *pgx.Conn, authService services.AuthService) (*models.User, error) {
	var username, password, email string
	var mobile int

	// Get user input for username
	color.New(color.FgCyan).Print("Enter username: ")
	fmt.Scan(&username)
	if !validation.IsValidUsername(username) {
		return nil, errors.New("invalid username: must be one word, alphanumeric, and can contain underscores")
	}

	// Get password input and validate
	password = utils.GetHiddenInput("Enter password: ")
	if !validation.IsValidPassword(password) {
		return nil, errors.New("invalid password: must be at least 8 characters, include an uppercase letter, a number, and a special character")
	}

	// Get email input and validate
	color.New(color.FgCyan).Print("Enter email: ")
	fmt.Scan(&email)
	if !validation.IsValidEmail(email) {
		return nil, errors.New("invalid email: must be a valid email address")
	}

	// Get mobile number input and validate
	color.New(color.FgCyan).Print("Enter mobile (10 digits): ")
	fmt.Scan(&mobile)
	if !validation.IsValidMobile(mobile) {
		return nil, errors.New("invalid mobile number: must be 10 digits")
	}

	// Hash the password before creating the user object
	hashedPassword := utils.HashPassword(password)

	// Create a new user object
	user := &models.User{
		Username: username,
		Password: hashedPassword,
		Email:    email,
		Mobile:   mobile,
		IsAdmin:  false,
		Role:     "user",
	}

	// Insert user into PostgreSQL database using the auth package
	err := authService.Signup(user)

	return user, err
}

// LoginUI handles user input and validation for PostgreSQL
func (ui *UI) LoginUI(authService services.AuthService) (*models.User, string, error) {
	var username, password string

	// Get user input for username
	color.New(color.FgCyan).Print("Enter username: ")
	fmt.Scan(&username)
	if username == "" {
		return nil, "", errors.New("username cannot be empty")
	}

	// Get user input for password
	password = utils.GetHiddenInput("Enter password: ")
	if password == "" {
		return nil, "", errors.New("password cannot be empty")
	}

	// Authenticate user with PostgreSQL using the auth package
	user, role, err := authService.Login(username, password)

	notifications, err := ui.notificationService.CheckNotification(username)
	if err != nil {
		color.New(color.FgRed).Println("Failed to check notifications:", err)
	} else if len(notifications) > 0 {
		color.New(color.FgGreen).Println("Notifications:")
		for _, notification := range notifications {
			color.New(color.FgGreen).Printf("Notification %d: %s\n", notification.Index, notification.Message)
		}
	} else {
		color.New(color.FgYellow).Println("No new notifications.")
	}

	return user, role, err
}
