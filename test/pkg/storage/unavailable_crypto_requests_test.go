package storage_test

import (
	"context"
	"cryptotracker/models"
	"cryptotracker/pkg/storage"
	//"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//const testConnString = "postgres://postgres:admin_password@localhost:5432/cryptotracker_test" // Test database

func setupTestDatabase(t *testing.T) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), testConnString)
	require.NoError(t, err, "failed to connect to test database")

	_, err = conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS unavailable_cryptos (
			id SERIAL PRIMARY KEY,
			crypto_symbol TEXT,
			username TEXT,
			request_message TEXT,
			status TEXT,
			timestamp TIMESTAMP
		)
	`)
	require.NoError(t, err, "failed to create test table")

	return conn
}

func teardownTestDatabase(t *testing.T, conn *pgx.Conn) {
	_, err := conn.Exec(context.Background(), `DELETE FROM unavailable_cryptos`)
	require.NoError(t, err, "failed to clean up test data")
	conn.Close(context.Background())
}

func TestPGUnavailableCryptoRequestRepository_SaveUnavailableCryptoRequest(t *testing.T) {
	conn := setupTestDatabase(t)
	defer teardownTestDatabase(t, conn)

	repo := storage.NewPGUnavailableCryptoRequestRepository(conn)

	tests := []struct {
		name    string
		request *models.UnavailableCryptoRequest
	}{
		{
			name: "valid unavailable crypto request",
			request: &models.UnavailableCryptoRequest{
				CryptoSymbol:   "ADA",
				UserName:       "testuser",
				RequestMessage: "Please add ADA to the tracker.",
				Status:         "pending",
				Timestamp:      time.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.SaveUnavailableCryptoRequest(conn, tt.request)
			assert.NoError(t, err)

			var count int
			query := `SELECT COUNT(*) FROM unavailable_cryptos WHERE crypto_symbol = $1 AND username = $2`
			err = conn.QueryRow(context.Background(), query, tt.request.CryptoSymbol, tt.request.UserName).Scan(&count)
			assert.NoError(t, err)
			assert.Equal(t, 1, count, "unavailable crypto request should be saved")
		})
	}
}
