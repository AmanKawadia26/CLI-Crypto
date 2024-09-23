package auth

import (
	"context"
	"cryptotracker/models"
	"cryptotracker/pkg/config"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
)

func Signup(conn *pgx.Conn, user *models.User) error {
	// Define columns for the SELECT query
	selectColumns := []string{"username"}
	selectCondition := "username = $1"

	// Use BuildSelectQuery to create a SELECT query for checking existing user
	selectQuery, err := config.BuildSelectQuery(selectColumns, "users", selectCondition)
	if err != nil {
		return fmt.Errorf("failed to build select query: %v", err)
	}

	// Check if the user already exists
	var existingUser models.User
	err = conn.QueryRow(context.Background(), selectQuery, user.Username).Scan(&existingUser.Username)
	if err == nil {
		return errors.New("user already exists")
	} else if err != pgx.ErrNoRows {
		log.Fatalf("failed to query user: %v", err)
		return err
	}

	// Define columns for the INSERT query
	insertColumns := []string{"userid", "username", "email", "mobile", "password", "role", "isadmin"}

	// Use BuildInsertQuery to create an INSERT query for user registration
	insertQuery, err := config.BuildInsertQuery("users", insertColumns)
	if err != nil {
		return fmt.Errorf("failed to build insert query: %v", err)
	}

	// Insert the new user
	_, err = conn.Exec(context.Background(), insertQuery, user.UserID, user.Username, user.Email, user.Mobile, user.Password, user.Role, user.IsAdmin)
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
		return err
	}

	return nil
}
