package notification

import (
	"context"
	"cryptotracker/internal/api"
	"cryptotracker/models"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/jackc/pgx/v4"
	"time"
)

// CheckNotification fetches and displays user notifications for PostgreSQL
func CheckNotification(conn *pgx.Conn, username string) {
	// Check Unavailable Crypto Requests notifications
	CheckUnavailableCryptoRequests(conn, username)

	// Check Price Alert Notifications
	CheckPriceAlerts(conn, username)
}

func CheckUnavailableCryptoRequests(conn *pgx.Conn, username string) {
	// Query the unavailable crypto requests with "Approved" or "Rejected" status for the logged-in user
	query := `SELECT crypto_symbol, status FROM unavailable_cryptos WHERE username = $1 AND status IN ('Approved', 'Rejected')`
	rows, err := conn.Query(context.Background(), query, username)
	if err != nil {
		color.Red("Error fetching unavailable crypto requests: %v", err)
		return
	}
	defer rows.Close()

	// Iterate through the user's requests
	for rows.Next() {
		var cryptoSymbol string
		var status string

		err := rows.Scan(&cryptoSymbol, &status)
		if err != nil {
			color.Red("Failed to scan unavailable crypto request: %v", err)
			continue
		}

		// Show notification based on the request's status
		switch status {
		case "Approved":
			color.Green("Admin has approved your request to add the cryptocurrency %s. You can now view it in the application.", cryptoSymbol)
		case "Rejected":
			color.Red("Your request to add the cryptocurrency %s has been rejected by the admin.", cryptoSymbol)
		}
	}

	if rows.Err() != nil {
		color.Red("Row iteration error while fetching unavailable crypto requests: %v", rows.Err())
	}
}

func CheckPriceAlerts(conn *pgx.Conn, username string) {
	// Query the price notifications for the logged-in user with "Pending" status
	query := `SELECT crypto_id, crypto, target_price FROM price_notifications WHERE username = $1 AND status = 'Pending'`
	rows, err := conn.Query(context.Background(), query, username)
	if err != nil {
		color.Red("Error fetching price notifications: %v", err)
		return
	}
	defer rows.Close()

	// Iterate through the user's price notifications
	for rows.Next() {
		var notification models.PriceNotification

		err := rows.Scan(&notification.CryptoID, &notification.Crypto, &notification.TargetPrice)
		if err != nil {
			color.Red("Failed to scan price notification: %v", err)
			continue
		}

		// Make an API call to check the current price of the cryptocurrency
		params := map[string]string{
			"id":      fmt.Sprintf("%d", notification.CryptoID),
			"convert": "USD",
		}
		response := api.GetAPIResponse("/quotes/latest", params)

		var result map[string]interface{}
		err = json.Unmarshal(response, &result)
		if err != nil {
			color.Red("Error unmarshalling API response: %v", err)
			continue
		}

		// Retrieve price information from the response
		data, dataOk := result["data"].(map[string]interface{})
		if !dataOk || data[fmt.Sprintf("%d", notification.CryptoID)] == nil {
			color.Red("Cryptocurrency data not found for ID: %d", notification.CryptoID)
			continue
		}

		cryptoData, ok := data[fmt.Sprintf("%d", notification.CryptoID)].(map[string]interface{})
		if !ok {
			color.Red("Unexpected data structure for crypto ID: %d", notification.CryptoID)
			continue
		}

		// Check if the price meets or exceeds the target price
		priceData, ok := cryptoData["quote"].(map[string]interface{})
		if !ok || priceData["USD"] == nil {
			color.Red("Price data not available for crypto ID: %d", notification.CryptoID)
			continue
		}

		currentPrice, ok := priceData["USD"].(map[string]interface{})["price"].(float64)
		if !ok {
			color.Red("Failed to retrieve current price for crypto ID: %d", notification.CryptoID)
			continue
		}

		// If the price meets the target price, update the status and notify the user
		if currentPrice >= notification.TargetPrice {
			color.Green("The cryptocurrency %s has reached your target price of $%.2f. Current price: $%.2f", notification.Crypto, notification.TargetPrice, currentPrice)
			notification.Status = "Served"
			notification.ServedAt = time.Now().Format(time.RFC3339)

			// Update the notification in the database
			_, err = conn.Exec(context.Background(),
				"UPDATE price_notifications SET status = 'Served', served_at = $1 WHERE crypto_id = $2 AND username = $3",
				notification.ServedAt, notification.CryptoID, username)
			if err != nil {
				color.Red("Error updating price notification status: %v", err)
			}
		}
	}

	if rows.Err() != nil {
		color.Red("Row iteration error while fetching price notifications: %v", rows.Err())
	}
}
