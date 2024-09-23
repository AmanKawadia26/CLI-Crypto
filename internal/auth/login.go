package auth

//
//import (
//	//"context"
//	"cryptotracker/internal/notification"
//	"cryptotracker/models"
//	//"cryptotracker/pkg/config"
//	"cryptotracker/pkg/utils"
//	"errors"
//	//"fmt"
//
//	"github.com/jackc/pgx/v4"
//)
//
//func Login(conn *pgx.Conn, username, password string) (*models.User, string, error) {
//
//	// Call the repository function to fetch the user from the database
//	user, err := LoginDBRepository(conn, username)
//	if err != nil {
//		return nil, "", err // Return the error if user is not found or any other DB issue
//	}
//
//	// Compare the provided password with the stored hashed password
//	hashedPassword := utils.HashPassword(password)
//	if user.Password != hashedPassword {
//		return nil, "", errors.New("invalid username or password")
//	}
//
//	// Check and display any notifications for the user
//	notification.CheckNotification(conn, username)
//
//	// Return the user object and role
//	return user, user.Role, nil
//}

//// LoginDBRepository handles the database interaction for user login.
//func LoginDBRepository(conn *pgx.Conn, username string) (*models.User, error) {
//
//	var user models.User
//
//	// Define the columns and conditions for the query
//	columns := []string{"userid", "username", "password", "email", "mobile", "role", "isadmin"}
//	condition := "username = $1"
//
//	// Build the query to fetch user information
//	query, err := config.BuildSelectQuery(columns, "users", condition)
//	if err != nil {
//		return nil, fmt.Errorf("failed to build query: %v", err)
//	}
//
//	// Execute the query to fetch the user information
//	err = conn.QueryRow(context.Background(), query, username).
//		Scan(&user.UserID, &user.Username, &user.Password, &user.Email, &user.Mobile, &user.Role, &user.IsAdmin)
//	if err != nil {
//		if err == pgx.ErrNoRows {
//			return nil, errors.New("user not found")
//		}
//		return nil, err // Return other DB-related errors
//	}
//
//	return &user, nil
//}
