package notification_test

import (
	"context"
	"cryptotracker/internal/notification"
	"cryptotracker/models"
	"cryptotracker/pkg/config"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
	"time"
)

var testConn *pgx.Conn

var (
	testDBURL = "postgres://postgres:admin_password@localhost:5432/cryptotracker_test" // Update this URL to match your test database setup
)

// Setup initializes the database connection and seeds the test database
func Setup() *pgx.Conn {
	// Connect to the test database
	conn, err := pgx.Connect(context.Background(), testDBURL)
	if err != nil {
		log.Fatalf("Failed to connect to the test database: %v", err)
	}

	// Clean up the tables
	_, err = conn.Exec(context.Background(), `
		TRUNCATE TABLE unavailable_cryptos, price_notifications RESTART IDENTITY;
	`)
	if err != nil {
		log.Fatalf("Failed to truncate tables: %v", err)
	}

	return conn
}

// TestMain sets up and tears down the test environment
func TestMain(m *testing.M) {
	testConn = Setup()
	defer testConn.Close(context.Background())

	code := m.Run()

	os.Exit(code)
}

func TestCheckNotification(t *testing.T) {
	conn := Setup()
	defer conn.Close(context.Background())

	currentTime := time.Now().Format(time.RFC3339)

	// Seed the database with test data for unavailable cryptos
	_, err := conn.Exec(context.Background(), `
		INSERT INTO unavailable_cryptos (username, crypto_symbol, status, request_message, timestamp) VALUES
		('testuser', 'abcd', 'Approved', 'Request to add abcd', $1),
		('testuser', 'efgh', 'Rejected', 'Request to add efgh', $2);
	`, currentTime, currentTime)
	require.NoError(t, err)

	// Seed the database with test data for price notifications
	_, err = conn.Exec(context.Background(), `
		INSERT INTO price_notifications (username, crypto_id, crypto, target_price, status, asked_at) VALUES
		('testuser', 1, 'BTC', 20000, 'Pending', $1),
		('testuser', 2, 'ETH', 1500, 'Pending', $2);
	`, time.Now(), time.Now())
	require.NoError(t, err)

	// Set up a test API key
	config.AppConfig.APIKey = "https://pro-api.coinmarketcap.com/v1"

	// Run the CheckNotification function
	notification.CheckNotification(conn, "testuser")

	var btcAlert models.PriceNotification
	err = conn.QueryRow(context.Background(), `
        SELECT crypto, target_price, status
        FROM price_notifications
        WHERE username = $1 AND crypto_id = $2
    `, "testuser", 1).Scan(&btcAlert.Crypto, &btcAlert.TargetPrice, &btcAlert.Status)
	require.NoError(t, err)
	assert.Equal(t, "Pending", btcAlert.Status)

	// Validate the ETH price alert (Pending)
	var ethAlert models.PriceNotification
	err = conn.QueryRow(context.Background(), `
        SELECT crypto, target_price, status
        FROM price_notifications
        WHERE username = $1 AND crypto_id = $2
    `, "testuser", 2).Scan(&ethAlert.Crypto, &ethAlert.TargetPrice, &ethAlert.Status)
	require.NoError(t, err)
	assert.Equal(t, "Pending", ethAlert.Status)
}

// TestCheckUnavailableCryptoRequests tests CheckUnavailableCryptoRequests function
func TestCheckUnavailableCryptoRequests(t *testing.T) {
	conn := Setup()
	defer conn.Close(context.Background())

	currentTime := time.Now().Format(time.RFC3339)

	// Seed the database with test data for unavailable cryptos
	_, err := conn.Exec(context.Background(), `
		INSERT INTO unavailable_cryptos (username, crypto_symbol, status, request_message, timestamp) VALUES
		('testuser', 'abcd', 'Approved', 'Request to add abcd', $1),
		('testuser', 'efgh', 'Rejected', 'Request to add efgh', $2);
	`, currentTime, currentTime)
	require.NoError(t, err)

	// Run the function
	notification.CheckUnavailableCryptoRequests(conn, "testuser")

	// Validate the results
	// Example: Capture output or validate the table's state after running the function.
}

func TestCheckPriceAlerts(t *testing.T) {
	conn := Setup()
	defer conn.Close(context.Background())

	// Seed the database
	_, err := conn.Exec(context.Background(), `
        INSERT INTO price_notifications (crypto_id, crypto, target_price, username, status, asked_at)
        VALUES (1, 'BTC', 20000, 'testuser', 'Pending', $1),
               (2, 'ETH', 1500, 'testuser', 'Pending', $2);
    `, time.Now(), time.Now())
	require.NoError(t, err)

	// Set up a test API key (mock or real)
	config.AppConfig.APIKey = "7deda392-59ec-47ff-a936-144f16086ed7"

	// Run the CheckPriceAlerts function
	notification.CheckPriceAlerts(conn, "testuser")

	// Validate the BTC price alert (Served)
	var btcAlert models.PriceNotification
	err = conn.QueryRow(context.Background(), `
        SELECT crypto, target_price, status
        FROM price_notifications
        WHERE username = $1 AND crypto_id = $2
    `, "testuser", 1).Scan(&btcAlert.Crypto, &btcAlert.TargetPrice, &btcAlert.Status)
	require.NoError(t, err)
	assert.Equal(t, "Pending", btcAlert.Status)

	// Validate the ETH price alert (Pending)
	var ethAlert models.PriceNotification
	err = conn.QueryRow(context.Background(), `
        SELECT crypto, target_price, status
        FROM price_notifications
        WHERE username = $1 AND crypto_id = $2
    `, "testuser", 2).Scan(&ethAlert.Crypto, &ethAlert.TargetPrice, &ethAlert.Status)
	require.NoError(t, err)
	assert.Equal(t, "Pending", ethAlert.Status)
}
