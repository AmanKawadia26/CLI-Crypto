////go:build !test
//// +build !test
//
//package ui
//
//import (
//	"cryptotracker/internal/auth"
//	"cryptotracker/models"
//	"cryptotracker/pkg/utils"
//	"cryptotracker/pkg/validation"
//	"errors"
//	"fmt"
//	"github.com/fatih/color"
//	"go.mongodb.org/mongo-driver/mongo"
//)
//
//// AuthenticateUser handles the login/signup process
//func AuthenticateUser(client *mongo.Client) (*models.User, string) {
//	for {
//		ClearScreen()
//		DisplayAuthMenu()
//
//		var choice int
//		color.New(color.FgCyan).Print("Enter your choice: ")
//		fmt.Scan(&choice)
//
//		switch choice {
//		case 1:
//			if user, Role, err := LoginUI(client); err == nil {
//				return user, Role
//			} else {
//				color.New(color.FgRed).Println("Login failed:", err)
//			}
//		case 2:
//			if _, err := SignupUI(client); err != nil {
//				color.New(color.FgRed).Println("Signup failed:", err)
//			} else {
//				color.New(color.FgGreen).Println("Signup successful. Please login.")
//			}
//		case 3:
//			color.New(color.FgYellow).Println("Exiting...")
//			return nil, ""
//		default:
//			color.New(color.FgRed).Println("Invalid choice, please try again.")
//		}
//	}
//}
//
//// SignupUI handles user input and validation
//func SignupUI(client *mongo.Client) (*models.User, error) {
//	var username, password, email string
//	var mobile int
//
//	// Get user input for username
//	color.New(color.FgCyan).Print("Enter username: ")
//	fmt.Scan(&username)
//	if !validation.IsValidUsername(username) {
//		return nil, errors.New("invalid username: must be one word, alphanumeric, and can contain underscores")
//	}
//
//	// Get password input and validate
//	password = utils.GetHiddenInput("Enter password: ")
//	if !validation.IsValidPassword(password) {
//		return nil, errors.New("invalid password: must be at least 8 characters, include an uppercase letter, a number, and a special character")
//	}
//
//	// Get email input and validate
//	color.New(color.FgCyan).Print("Enter email: ")
//	fmt.Scan(&email)
//	if !validation.IsValidEmail(email) {
//		return nil, errors.New("invalid email: must be a valid email address")
//	}
//
//	// Get mobile number input and validate
//	color.New(color.FgCyan).Print("Enter mobile (10 digits): ")
//	fmt.Scan(&mobile)
//	if !validation.IsValidMobile(mobile) {
//		return nil, errors.New("invalid mobile number: must be 10 digits")
//	}
//
//	// Hash the password before creating the user object
//	hashedPassword := utils.HashPassword(password)
//
//	// Create and return a new user object
//	user := &models.User{
//		Username: username,
//		Password: hashedPassword,
//		Email:    email,
//		Mobile:   mobile,
//		IsAdmin:  false,
//		Role:     "user",
//	}
//
//	err := auth.Signup(client, user)
//
//	return user, err
//}
//
//// LoginUI handles user input and validation
//func LoginUI(client *mongo.Client) (*models.User, string, error) {
//	var username, password string
//
//	// Get user input for username
//	color.New(color.FgCyan).Print("Enter username: ")
//	fmt.Scan(&username)
//	if username == "" {
//		return nil, "", errors.New("username cannot be empty")
//	}
//
//	// Get user input for password
//	password = utils.GetHiddenInput("Enter password: ")
//	if password == "" {
//		return nil, "", errors.New("password cannot be empty")
//	}
//
//	user, role, err := auth.Login(client, username, password)
//
//	return user, role, err
//}

//go:build !test
// +build !test

package ui

import (
	//"context"
	"cryptotracker/internal/auth"
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
func AuthenticateUser(conn *pgx.Conn) (*models.User, string) {
	for {
		ClearScreen()
		DisplayAuthMenu()

		var choice int
		color.New(color.FgCyan).Print("Enter your choice: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			if user, Role, err := LoginUI(conn); err == nil {
				return user, Role
			} else {
				color.New(color.FgRed).Println("Login failed:", err)
			}
		case 2:
			if _, err := SignupUI(conn); err != nil {
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
func SignupUI(conn *pgx.Conn) (*models.User, error) {
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
		UserID:   "0",
		Username: username,
		Password: hashedPassword,
		Email:    email,
		Mobile:   mobile,
		IsAdmin:  false,
		Role:     "user",
	}

	// Insert user into PostgreSQL database using the auth package
	err := auth.Signup(conn, user)

	return user, err
}

// LoginUI handles user input and validation for PostgreSQL
func LoginUI(conn *pgx.Conn) (*models.User, string, error) {
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
	user, role, err := auth.Login(conn, username, password)

	return user, role, err
}
