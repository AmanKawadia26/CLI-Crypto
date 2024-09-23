package storage

import (
	"context"
	"cryptotracker/models"
	"fmt"
	"github.com/jackc/pgx/v4"
	"time"
)

// UnavailableCryptoRequestRepository defines methods to interact with unavailable crypto requests storage.
type UnavailableCryptoRequestRepository interface {
	SaveUnavailableCryptoRequest(request *models.UnavailableCryptoRequest) error
}

// PGUnavailableCryptoRequestRepository is the PostgreSQL implementation of UnavailableCryptoRequestRepository.
type PGUnavailableCryptoRequestRepository struct {
	conn *pgx.Conn
}

// NewPGUnavailableCryptoRequestRepository creates a new instance of PGUnavailableCryptoRequestRepository.
func NewPGUnavailableCryptoRequestRepository(conn *pgx.Conn) *PGUnavailableCryptoRequestRepository {
	return &PGUnavailableCryptoRequestRepository{conn: conn}
}

const (
	unavailableCryptoTable = "unavailable_cryptos"
)

// SaveUnavailableCryptoRequest saves a new unavailable crypto request in PostgreSQL
func (r *PGUnavailableCryptoRequestRepository) SaveUnavailableCryptoRequest(conn *pgx.Conn, request *models.UnavailableCryptoRequest) error {
	query := `
		INSERT INTO unavailable_cryptos (crypto_symbol, username, request_message, status, timestamp)
		VALUES ($1, $2, $3, $4, $5)
	`

	// Set the request creation time if needed
	request.Timestamp = time.Now()

	_, err := conn.Exec(context.Background(), query, request.CryptoSymbol, request.UserName, request.RequestMessage, request.Status, request.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to save unavailable crypto request: %v", err)
	}

	return nil
}

//// GetAllUnavailableCryptoRequests retrieves all unavailable crypto requests from PostgreSQL
//func GetAllUnavailableCryptoRequests(conn *pgx.Conn) ([]*models.UnavailableCryptoRequest, error) {
//	query := `
//		SELECT crypto_symbol, username, request_message, status, timestamp
//		FROM unavailable_cryptos
//	`
//
//	rows, err := conn.Query(context.Background(), query)
//	if err != nil {
//		return nil, fmt.Errorf("error fetching unavailable crypto requests: %v", err)
//	}
//	defer rows.Close() // Ensure rows are closed after processing
//
//	var requests []*models.UnavailableCryptoRequest
//
//	for rows.Next() {
//		var request models.UnavailableCryptoRequest
//		if err := rows.Scan(&request.CryptoSymbol, &request.UserName, &request.RequestMessage, &request.Status, &request.Timestamp); err != nil {
//			return nil, fmt.Errorf("failed to scan unavailable crypto request: %v", err)
//		}
//		requests = append(requests, &request)
//	}
//
//	if err := rows.Err(); err != nil {
//		return nil, fmt.Errorf("rows error while fetching unavailable crypto requests: %v", err)
//	}
//
//	return requests, nil
//}
