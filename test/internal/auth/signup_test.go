package auth

import (
	"context"
	"cryptotracker/internal/auth"
	"cryptotracker/models"
	"errors"
	//"github.com/jackc/pgx/v4"
	"testing"
)

// TestSignup function for PostgreSQL
func TestSignup(t *testing.T) {
	conn, err := setupPostgres()
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer conn.Close(context.Background())

	tests := []struct {
		name          string
		existingUsers []*models.User
		newUser       *models.User
		expectedErr   error
	}{
		{
			name:          "Successful signup",
			existingUsers: []*models.User{}, // No existing users
			newUser: &models.User{
				Username: "test_user_success",
				Password: "hashed_password",
				Email:    "success@example.com",
				Mobile:   1234567890,
				Role:     "user",
				IsAdmin:  false,
			},
			expectedErr: nil,
		},
		{
			name: "User already exists",
			existingUsers: []*models.User{
				{Username: "existing_user", Password: "hashed_password", Email: "existing@example.com", Mobile: 1234567890},
			},
			newUser: &models.User{
				Username: "existing_user", // Same username as an existing user
				Password: "hashed_password",
				Email:    "new@example.com",
				Mobile:   1234567890,
				Role:     "user",
				IsAdmin:  false,
			},
			expectedErr: errors.New("user already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Insert existing users if any
			for _, u := range tt.existingUsers {
				err := insertTestUser(conn, u)
				if err != nil {
					t.Fatalf("Failed to insert existing user: %v", err)
				}
			}

			// Prepare for test by inserting a user
			testUsername := tt.newUser.Username
			// Perform Signup
			err := auth.Signup(conn, tt.newUser)

			if (err != nil) != (tt.expectedErr != nil) || (err != nil && err.Error() != tt.expectedErr.Error()) {
				t.Errorf("Signup() error = %v, expectedErr %v", err, tt.expectedErr)
				// Cleanup test data
				if err == nil {
					cleanupTestUser(conn, testUsername)
				}
				return
			}

			// If signup was successful, check if the new user was saved
			if tt.expectedErr == nil {
				var count int
				err := conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM users WHERE username=$1", tt.newUser.Username).Scan(&count)
				if err != nil {
					t.Fatalf("Failed to get user count from PostgreSQL: %v", err)
				}
				if count != 1 {
					t.Errorf("Signup() failed to save the user %v", tt.newUser.Username)
				}
			}

			// Cleanup test data
			if tt.expectedErr == nil {
				cleanupTestUser(conn, testUsername)
			}
		})
	}
}
