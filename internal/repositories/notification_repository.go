package repositories

import (
	"context"
	"cryptotracker/models"
	"cryptotracker/pkg/config"
	"fmt"
	"github.com/jackc/pgx/v4"
	"time"
)

type NotificationRepository interface {
	CheckUnavailableCryptoRequestsRepo(username string) (pgx.Rows, error)
	CheckPriceAlertsRepo(username string) ([]models.PriceNotification, error)
	UpdatePriceNotificationStatusRepo(notification *models.PriceNotification, username string, currentPrice float64) error
}

type PostgresNotificationRepository struct {
	conn *pgx.Conn
}

func NewPostgresNotificationRepository(conn *pgx.Conn) *PostgresNotificationRepository {
	return &PostgresNotificationRepository{conn: conn}
}

func (r *PostgresNotificationRepository) CheckUnavailableCryptoRequestsRepo(username string) (pgx.Rows, error) {
	columns := []string{"crypto_symbol", "status"}
	condition := "username = $1 AND status IN ('Approved', 'Rejected')"

	query, err := config.BuildSelectQuery(columns, "unavailable_cryptos", condition)
	if err != nil {
		return nil, err
	}

	rows, err := r.conn.Query(context.Background(), query, username)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *PostgresNotificationRepository) CheckPriceAlertsRepo(username string) ([]models.PriceNotification, error) {
	columns := []string{"crypto_id", "crypto", "target_price"}
	condition := "username = $1 AND status = 'Pending'"

	query, err := config.BuildSelectQuery(columns, "price_notifications", condition)
	if err != nil {
		return nil, err
	}

	rows, err := r.conn.Query(context.Background(), query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.PriceNotification
	for rows.Next() {
		var notification models.PriceNotification
		err := rows.Scan(&notification.CryptoID, &notification.Crypto, &notification.TargetPrice)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return notifications, nil
}

func (r *PostgresNotificationRepository) UpdatePriceNotificationStatusRepo(notification *models.PriceNotification, username string, cryptoPrice float64) error {
	// Only update if the crypto price meets or exceeds the target price
	if cryptoPrice >= notification.TargetPrice {
		// Update the status and served time
		notification.Status = "Served"
		notification.ServedAt = time.Now().Format(time.RFC3339)

		// Build the update query
		updateColumns := []string{"status", "served_at"}
		// Add condition to check for the correct notification using crypto_id, username, and price
		updateCondition := "crypto_id = $3 AND username = $4 AND target_price <= $5 AND status = 'Pending'"

		// Generate the query
		updateQuery, err := config.BuildUpdateQuery("price_notifications", updateColumns, updateCondition)
		if err != nil {
			return err
		}

		// Set the context for the query
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Execute the query, passing the current values and crypto price
		_, err = r.conn.Exec(ctx, updateQuery, notification.Status, notification.ServedAt, notification.CryptoID, username, notification.TargetPrice)
		if err != nil {
			return err
		}
	}

	return nil
}

// SavePriceNotification saves a single notification to PostgreSQL
func (r *PostgresNotificationRepository) SavePriceNotification(conn *pgx.Conn, notification *models.PriceNotification) error {
	columns := []string{"crypto_id", "crypto", "target_price", "username", "asked_at", "status", "served_at"}
	query, err := config.BuildInsertQuery("price_notifications", columns)
	if err != nil {
		return fmt.Errorf("failed to build insert query: %v", err)
	}

	// Handle the case where ServedAt is not set (i.e., it's pending)
	var servedAt interface{}
	if notification.ServedAt == "" {
		servedAt = nil // Use nil to insert NULL into the database
	} else {
		servedAt = notification.ServedAt
	}

	_, err = conn.Exec(context.Background(), query,
		notification.CryptoID,
		notification.Crypto,
		notification.TargetPrice,
		notification.Username,
		notification.AskedAt,
		notification.Status,
		servedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save price notification: %v", err)
	}

	return nil
}
