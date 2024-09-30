package repositories

import (
	"context"
	"cryptotracker/models"
	"cryptotracker/pkg/config"
	"fmt"
	"github.com/jackc/pgx/v4"
)

type AdminRepository interface {
	ChangeUserStatus(username string) error
	DeleteUser(username string) error
	ViewUserProfiles() ([]*models.User, error)
	ManageUserRequests() ([]*models.UnavailableCryptoRequest, error)
	ManageSpecificCryptoRequests(cryptoSymbol string) ([]*models.UnavailableCryptoRequest, error)
	SaveUnavailableCryptoRequest(requests []*models.UnavailableCryptoRequest) error
}

type PostgresAdminRepository struct {
	conn *pgx.Conn
}

func NewPostgresAdminRepository(conn *pgx.Conn) *PostgresAdminRepository {
	return &PostgresAdminRepository{conn: conn}
}

func (r *PostgresAdminRepository) ChangeUserStatus(username string) error {
	query, err := config.BuildSelectQuery([]string{"username", "role"}, "users", "username=$1")
	if err != nil {
		return fmt.Errorf("failed to build query: %v", err)
	}

	row := r.conn.QueryRow(context.Background(), query, username)

	var user models.User
	err = row.Scan(&user.Username, &user.Role)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to query user: %v", err)
	}

	if user.Role == "admin" {
		return fmt.Errorf("user is already an admin")
	}

	updateQuery, err := config.BuildUpdateQuery("users", []string{"role", "isadmin"}, "username=$3")
	if err != nil {
		return fmt.Errorf("failed to build update query: %v", err)
	}

	_, err = r.conn.Exec(context.Background(), updateQuery, "admin", true, username)
	if err != nil {
		return fmt.Errorf("failed to update user role: %v", err)
	}

	return nil
}

func (r *PostgresAdminRepository) DeleteUser(username string) error {
	tx, err := r.conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			tx.Commit(context.Background())
		}
	}()

	deleteUserQuery, err := config.BuildDeleteQuery("users", "username=$1")
	if err != nil {
		return fmt.Errorf("failed to build delete user query: %v", err)
	}
	_, err = tx.Exec(context.Background(), deleteUserQuery, username)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	// Delete related price notifications and unavailable crypto requests
	_, err = tx.Exec(context.Background(), "DELETE FROM price_notifications WHERE username=$1", username)
	if err != nil {
		return fmt.Errorf("failed to delete user price notifications: %v", err)
	}

	_, err = tx.Exec(context.Background(), "DELETE FROM unavailable_cryptos WHERE username=$1", username)
	if err != nil {
		return fmt.Errorf("failed to delete user unavailable crypto requests: %v", err)
	}

	return nil
}

func (r *PostgresAdminRepository) ViewUserProfiles() ([]*models.User, error) {
	query, err := config.BuildSelectQuery([]string{"userid", "username", "email", "mobile", "isadmin", "role"}, "users", "")
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	rows, err := r.conn.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %v", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.Email, &user.Mobile, &user.IsAdmin, &user.Role); err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return users, nil
}

func (r *PostgresAdminRepository) ManageUserRequests() ([]*models.UnavailableCryptoRequest, error) {
	query, err := config.BuildSelectQuery([]string{"crypto_symbol", "username", "request_message", "status", "timestamp"}, "unavailable_cryptos", "")
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	rows, err := r.conn.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to query unavailable crypto requests: %v", err)
	}
	defer rows.Close()

	var requests []*models.UnavailableCryptoRequest
	for rows.Next() {
		var request models.UnavailableCryptoRequest
		if err := rows.Scan(&request.CryptoSymbol, &request.UserName, &request.RequestMessage, &request.Status, &request.Timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan request: %v", err)
		}
		requests = append(requests, &request)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return requests, nil
}

func (r *PostgresAdminRepository) ManageSpecificCryptoRequests(cryptoSymbol string) ([]*models.UnavailableCryptoRequest, error) {
	condition := "crypto_symbol = $1"
	query, err := config.BuildSelectQuery([]string{"crypto_symbol", "username", "request_message", "status", "timestamp"}, "unavailable_cryptos", condition)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	rows, err := r.conn.Query(context.Background(), query, cryptoSymbol)
	if err != nil {
		return nil, fmt.Errorf("failed to query unavailable crypto requests: %v", err)
	}
	defer rows.Close()

	var requests []*models.UnavailableCryptoRequest
	for rows.Next() {
		var request models.UnavailableCryptoRequest
		if err := rows.Scan(&request.CryptoSymbol, &request.UserName, &request.RequestMessage, &request.Status, &request.Timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan request: %v", err)
		}
		requests = append(requests, &request)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return requests, nil
}

func (r *PostgresAdminRepository) SaveUnavailableCryptoRequest(requests []*models.UnavailableCryptoRequest) error {

	var err error
	var query string

	for _, request := range requests {
		query, err = config.BuildInsertQuery("unavailable_cryptos", []string{"crypto_symbol", "username", "request_message", "status", "timestamp"})
		if err != nil {
			return fmt.Errorf("failed to build insert query: %v", err)
		}

		_, err = r.conn.Exec(context.Background(), query, request.CryptoSymbol, request.UserName, request.RequestMessage, request.Status, request.Timestamp)
		if err != nil {
			return fmt.Errorf("failed to insert request: %v", err)
		}

		return nil
	}
	return nil
}
