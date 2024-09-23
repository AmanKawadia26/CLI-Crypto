package auth

import (
	"context"
	"cryptotracker/internal/notification"
	"cryptotracker/models"
	"cryptotracker/pkg/utils"
	"errors"
	"github.com/jackc/pgx/v4"
)

// Login handles the login process and verifies user credentials in PostgreSQL
func Login(conn *pgx.Conn, username, password string) (*models.User, string, error) {
	var user models.User

	// Fetch the user from PostgreSQL using the username
	err := conn.QueryRow(context.Background(),
		"SELECT userid, username, password, email, mobile, role, isadmin FROM users WHERE username=$1", username).
		Scan(&user.UserID, &user.Username, &user.Password, &user.Email, &user.Mobile, &user.Role, &user.IsAdmin)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, "", errors.New("user not found")
		}
		return nil, "", err
	}

	// Verify the hashed password
	hashedPassword := utils.HashPassword(password)
	if user.Password != hashedPassword {
		return nil, "", errors.New("invalid username or password")
	}

	// Check and display notifications for the user
	notification.CheckNotification(conn, username)

	return &user, user.Role, nil
}
