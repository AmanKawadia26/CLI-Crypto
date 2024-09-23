package auth

import (
	"context"
	"cryptotracker/models"
	"errors"
	"github.com/jackc/pgx/v4"
	"log"
)

// Signup handles the signup process using PostgreSQL
func Signup(conn *pgx.Conn, user *models.User) error {
	// Check if the user already exists in PostgreSQL
	var existingUser models.User
	err := conn.QueryRow(context.Background(), "SELECT username FROM users WHERE username=$1", user.Username).Scan(&existingUser.Username)
	if err == nil {
		return errors.New("user already exists")
	} else if err != pgx.ErrNoRows {
		log.Fatalf("failed to query user: %v", err)
		return err
	}

	// Insert the new user into PostgreSQL
	_, err = conn.Exec(context.Background(),
		"INSERT INTO users (userid, username, email, mobile, password, role, isadmin) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		user.UserID, user.Username, user.Email, user.Mobile, user.Password, user.Role, user.IsAdmin)
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
		return err
	}

	return nil
}
