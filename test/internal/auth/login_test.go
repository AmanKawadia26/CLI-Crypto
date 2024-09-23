package auth

import (
	"context"
	"cryptotracker/internal/auth"
	"cryptotracker/models"
	"cryptotracker/pkg/utils"
	"errors"
	"github.com/jackc/pgx/v4"
	"testing"
)

// Setup PostgreSQL connection
func setupPostgres() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:admin_password@localhost:5432/cryptotracker")
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Cleanup test data by username from PostgreSQL
func cleanupTestUser(conn *pgx.Conn, username string) {
	conn.Exec(context.Background(), "DELETE FROM users WHERE username=$1", username)
}

// Insert test user into PostgreSQL
func insertTestUser(conn *pgx.Conn, user *models.User) error {
	_, err := conn.Exec(context.Background(),
		"INSERT INTO users (username, password, email, mobile, role, isadmin) VALUES ($1, $2, $3, $4, $5, $6)",
		user.Username, user.Password, user.Email, user.Mobile, user.Role, user.IsAdmin)
	return err
}

// TestLogin function for PostgreSQL
func TestLogin(t *testing.T) {
	conn, err := setupPostgres()
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer conn.Close(context.Background())

	tests := []struct {
		name          string
		existingUser  *models.User
		loginUsername string
		loginPassword string
		expectedUser  *models.User
		expectedRole  string
		expectedErr   error
	}{
		{
			name: "Invalid username",
			existingUser: &models.User{
				UserID:   "1",
				Username: "test_user_invalid_username",
				Password: utils.HashPassword("password123"),
				Email:    "invalid@example.com",
				Mobile:   1234567890,
				Role:     "user",
				IsAdmin:  false,
			},
			loginUsername: "nonexistent_user",
			loginPassword: "password123",
			expectedUser:  nil,
			expectedRole:  "",
			expectedErr:   errors.New("user not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Insert existing user if any
			if tt.existingUser != nil {
				err := insertTestUser(conn, tt.existingUser)
				if err != nil {
					t.Fatalf("Failed to insert existing user: %v", err)
				}
			}

			// Perform Login
			user, role, err := auth.Login(conn, tt.loginUsername, tt.loginPassword)

			// Check for expected error
			if (err != nil) != (tt.expectedErr != nil) || (err != nil && err.Error() != tt.expectedErr.Error()) {
				t.Errorf("Login() error = %v, expectedErr %v", err, tt.expectedErr)
				// Cleanup test data
				if tt.existingUser != nil {
					cleanupTestUser(conn, tt.existingUser.Username)
				}
				return
			}

			// Check if user and role match the expected values
			if user != nil {
				if user.Username != tt.expectedUser.Username || user.Password != tt.expectedUser.Password ||
					user.Email != tt.expectedUser.Email || user.Mobile != tt.expectedUser.Mobile ||
					user.UserID != tt.expectedUser.UserID {
					t.Errorf("Login() user = %v, expectedUser %v", user, tt.expectedUser)
				}
				if role != tt.expectedRole {
					t.Errorf("Login() role = %v, expectedRole %v", role, tt.expectedRole)
				}
			} else {
				if tt.expectedUser != nil {
					t.Errorf("Login() expected user %v, but got nil", tt.expectedUser)
				}
			}

			// Cleanup test data
			if tt.existingUser != nil {
				cleanupTestUser(conn, tt.existingUser.Username)
			}
		})
	}
}
