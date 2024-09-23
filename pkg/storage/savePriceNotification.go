package storage

import (
	"context"
	"cryptotracker/models"
	"fmt"
	"github.com/jackc/pgx/v4"
)

const notificationsTable = "price_notifications"

// NotificationRepository defines the methods for saving and loading notifications.
type NotificationRepository interface {
	SavePriceNotification(conn *pgx.Conn, notification *models.PriceNotification) error
	//LoadPriceNotifications() ([]*models.PriceNotification, error)
}

// PGNotificationRepository is the PostgreSQL implementation of NotificationRepository.
type PGNotificationRepository struct {
	conn *pgx.Conn
}

// NewPGNotificationRepository creates a new instance of PGNotificationRepository.
func NewPGNotificationRepository(conn *pgx.Conn) *PGNotificationRepository {
	return &PGNotificationRepository{conn: conn}
}

// SavePriceNotification saves a single notification to PostgreSQL
func (r *PGNotificationRepository) SavePriceNotification(conn *pgx.Conn, notification *models.PriceNotification) error {
	query := `INSERT INTO price_notifications (crypto_id, crypto, target_price, username, asked_at, status, served_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	// Handle the case where ServedAt is not set (i.e., it's pending)
	var servedAt interface{}
	if notification.ServedAt == "" {
		servedAt = nil // Use nil to insert NULL into the database
	} else {
		servedAt = notification.ServedAt
	}

	_, err := conn.Exec(context.Background(), query,
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
