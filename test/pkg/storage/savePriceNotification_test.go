package storage_test

import (
	"context"
	"cryptotracker/models"
	"cryptotracker/pkg/storage"
	//"fmt"
	//"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testConnString = "postgres://postgres:admin_password@localhost:5432/cryptotracker_test" // Update with your DB credentials

// TestPGNotificationRepository_SavePriceNotification tests the SavePriceNotification method.
func TestPGNotificationRepository_SavePriceNotification(t *testing.T) {
	conn := setupTestDatabase(t)
	defer teardownTestDatabase(t, conn)

	repo := storage.NewPGNotificationRepository(conn)

	tests := []struct {
		name         string
		notification *models.PriceNotification
		expectedErr  bool
	}{
		{
			name: "valid notification with served_at",
			notification: &models.PriceNotification{
				CryptoID:    1,
				Crypto:      "Bitcoin",
				TargetPrice: 50000.0,
				Username:    "testuser",
				AskedAt:     time.Now().Format(time.RFC3339),
				Status:      "pending",
				ServedAt:    time.Now().Format(time.RFC3339),
			},
			expectedErr: false,
		},
		{
			name: "valid notification without served_at",
			notification: &models.PriceNotification{
				CryptoID:    2,
				Crypto:      "Ethereum",
				TargetPrice: 3000.0,
				Username:    "testuser2",
				AskedAt:     time.Now().Format(time.RFC3339),
				Status:      "pending",
				ServedAt:    "",
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.SavePriceNotification(conn, tt.notification)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Optionally, verify the notification was inserted
				var count int
				query := `SELECT COUNT(*) FROM price_notifications WHERE crypto_id = $1 AND username = $2`
				err = conn.QueryRow(context.Background(), query, tt.notification.CryptoID, tt.notification.Username).Scan(&count)
				assert.NoError(t, err)
				assert.Equal(t, 1, count, "notification should be saved")
			}
		})
	}
}
