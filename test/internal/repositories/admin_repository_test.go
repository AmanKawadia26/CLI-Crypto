package repositories_test

import (
	"context"
	"cryptotracker/internal/repositories"
	"fmt"
	"testing"
	"time"

	"cryptotracker/models"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// Adjust these constants to match your test database configuration
	testDBURL = "postgres://postgres:admin_password@localhost:5432/cryptotracker_test?sslmode=disable"
)

var (
	conn *pgx.Conn
)

func setupTestDB() (*pgx.Conn, func(), error) {
	// Connect to the test database
	conn, err := pgx.Connect(context.Background(), testDBURL)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	// Create a cleanup function
	cleanup := func() {
		conn.Close(context.Background())
	}

	// Setup test schema and data
	_, err = conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS users (
			userid TEXT NOT NULL,
			username TEXT NOT NULL,
			email TEXT NOT NULL,
			mobile BIGINT NOT NULL,
			isadmin BOOLEAN NOT NULL,
			role TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS price_notifications (
			id SERIAL PRIMARY KEY,
			username TEXT NOT NULL,
			crypto_symbol TEXT NOT NULL,
			target_price DECIMAL NOT NULL,
			timestamp TIMESTAMPTZ NOT NULL
		);
		CREATE TABLE IF NOT EXISTS unavailable_cryptos (
			id SERIAL PRIMARY KEY,
			crypto_symbol TEXT NOT NULL,
			username TEXT NOT NULL,
			request_message TEXT NOT NULL,
			status TEXT NOT NULL,
			timestamp TIMESTAMPTZ NOT NULL
		);
	`)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to set up database schema: %v", err)
	}

	return conn, cleanup, nil
}

func TestChangeUserStatus(t *testing.T) {
	conn, cleanup, err := setupTestDB()
	require.NoError(t, err)
	defer cleanup()

	repo := repositories.NewPostgresAdminRepository(conn)

	// Insert a user with a non-admin role
	_, err = conn.Exec(context.Background(), "INSERT INTO users (userid, username, email, mobile, isadmin, role) VALUES ($1, $2, $3, $4, $5, $6)",
		"0", "testuser", "testuser@example.com", "1234567890", false, "user")
	require.NoError(t, err)

	// Change the user's role to admin
	err = repo.ChangeUserStatus(conn, "testuser")
	require.NoError(t, err)

	// Verify the user role has been updated to admin
	row := conn.QueryRow(context.Background(), "SELECT role FROM users WHERE username=$1", "testuser")
	var role string
	err = row.Scan(&role)
	require.NoError(t, err)
	assert.Equal(t, "admin", role)

	// Test failure case: user is already an admin
	err = repo.ChangeUserStatus(conn, "testuser")
	assert.EqualError(t, err, "user is already an admin")

	// Test failure case: user not found
	err = repo.ChangeUserStatus(conn, "nonexistentuser")
	assert.EqualError(t, err, "user not found")
}

func TestDeleteUser(t *testing.T) {
	conn, cleanup, err := setupTestDB()
	require.NoError(t, err)
	defer cleanup()

	repo := repositories.NewPostgresAdminRepository(conn)

	_, err = conn.Exec(context.Background(), "INSERT INTO users (userid, username, email, mobile, isadmin, role) VALUES ($1, $2, $3, $4, $5, $6)",
		"0", "testuser", "testuser@example.com", "1234567890", false, "user")
	require.NoError(t, err)

	err = repo.DeleteUser(conn, "testuser")
	assert.NoError(t, err)

	// Test failure case: transaction error
	// This would require further configuration, such as simulating transaction errors
}

func TestViewUserProfiles(t *testing.T) {
	conn, cleanup, err := setupTestDB()
	require.NoError(t, err)
	defer cleanup()

	repo := repositories.NewPostgresAdminRepository(conn)

	// Insert a user with admin role for completeness
	_, err = conn.Exec(context.Background(), "INSERT INTO users (userid, username, email, mobile, isadmin, role) VALUES ($1, $2, $3, $4, $5, $6)",
		"0", "adminuser", "adminuser@example.com", 1234567890, true, "admin")
	require.NoError(t, err)

	// Insert a user with non-admin role for testing
	_, err = conn.Exec(context.Background(), "INSERT INTO users (userid, username, email, mobile, isadmin, role) VALUES ($1, $2, $3, $4, $5, $6)",
		"0", "testuser", "testuser@example.com", 987654321, false, "user")
	require.NoError(t, err)

	// Fetch user profiles
	profiles, err := repo.ViewUserProfiles(conn)
	if err != nil {
		t.Fatalf("Failed to fetch user profiles: %v", err)
	}

	// Ensure there is at least one profile
	assert.NotEmpty(t, profiles, "Expected at least one user profile")

	// Verify the details of the fetched profile(s)
	if len(profiles) > 0 {
		// Assert that the inserted test user is among the fetched profiles
		var found bool
		for _, profile := range profiles {
			if profile.Username == "testuser" {
				found = true
				assert.Equal(t, "testuser@example.com", profile.Email)
				assert.Equal(t, 987654321, profile.Mobile)
				assert.False(t, profile.IsAdmin)
				assert.Equal(t, "user", profile.Role)
				break
			}
		}
		assert.True(t, found, "Expected testuser to be in the list of profiles")
	}
}

func TestManageUserRequests(t *testing.T) {
	conn, cleanup, err := setupTestDB()
	require.NoError(t, err)
	defer cleanup()

	repo := repositories.NewPostgresAdminRepository(conn)

	_, err = conn.Exec(context.Background(), "INSERT INTO unavailable_cryptos (crypto_symbol, username, request_message, status, timestamp) VALUES ($1, $2, $3, $4, $5)",
		"BTC", "testuser", "Request message", "pending", time.Now())
	require.NoError(t, err)

	requests, err := repo.ManageUserRequests(conn)
	assert.NoError(t, err)
	assert.Len(t, requests, 1)
	assert.Equal(t, "BTC", requests[0].CryptoSymbol)
}

func TestSaveUnavailableCryptoRequest(t *testing.T) {
	conn, cleanup, err := setupTestDB()
	require.NoError(t, err)
	defer cleanup()

	repo := repositories.NewPostgresAdminRepository(conn)

	request := &models.UnavailableCryptoRequest{
		CryptoSymbol:   "BTC",
		UserName:       "testuser",
		RequestMessage: "Request message",
		Status:         "pending",
		Timestamp:      time.Now(),
	}

	err = repo.SaveUnavailableCryptoRequest(conn, request)
	assert.NoError(t, err)

	// Test failure case: simulate insert error if needed
}
