package repositories

import (
	"context"
	"cryptotracker/models"
	"fmt"
	"github.com/jackc/pgx/v4"
	//"time"
)

type AdminRepository interface {
	ChangeUserStatus(conn *pgx.Conn, username string) error
	DeleteUser(conn *pgx.Conn, username string) error
	ViewUserProfiles(conn *pgx.Conn) ([]*models.User, error)
	ManageUserRequests(conn *pgx.Conn) ([]*models.UnavailableCryptoRequest, error)
	SaveUnavailableCryptoRequest(conn *pgx.Conn, request *models.UnavailableCryptoRequest) error
}

type PostgresAdminRepository struct {
	conn *pgx.Conn
}

func NewPostgresAdminRepository(conn *pgx.Conn) *PostgresAdminRepository {
	return &PostgresAdminRepository{conn: conn}
}

func (r *PostgresAdminRepository) ChangeUserStatus(conn *pgx.Conn, username string) error {
	// Load user
	row := conn.QueryRow(context.Background(), "SELECT username, role FROM users WHERE username=$1", username)

	var user models.User
	err := row.Scan(&user.Username, &user.Role)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to query user: %v", err)
	}

	if user.Role == "admin" {
		return fmt.Errorf("user is already an admin")
	}

	// Update user to admin
	_, err = conn.Exec(context.Background(), "UPDATE users SET role=$1, isadmin=$2 WHERE username=$3", "admin", true, username)
	if err != nil {
		return fmt.Errorf("failed to update user role: %v", err)
	}

	return nil
}

func (r *PostgresAdminRepository) DeleteUser(conn *pgx.Conn, username string) error {
	// Begin a transaction to ensure atomicity
	tx, err := conn.Begin(context.Background())
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

	// Delete user from users table
	_, err = tx.Exec(context.Background(), "DELETE FROM users WHERE username=$1", username)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	// Delete related price notifications from price_notifications table
	_, err = tx.Exec(context.Background(), "DELETE FROM price_notifications WHERE username=$1", username)
	if err != nil {
		return fmt.Errorf("failed to delete user price notifications: %v", err)
	}

	// Delete related unavailable crypto requests from unavailable_cryptos table
	_, err = tx.Exec(context.Background(), "DELETE FROM unavailable_cryptos WHERE username=$1", username)
	if err != nil {
		return fmt.Errorf("failed to delete user unavailable crypto requests: %v", err)
	}

	// If all deletions succeeded, commit the transaction
	return nil
}

func (r *PostgresAdminRepository) ViewUserProfiles(conn *pgx.Conn) ([]*models.User, error) {
	var users []*models.User

	// Get all users
	rows, err := conn.Query(context.Background(), "SELECT userid, username, email, mobile, isadmin, role FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %v", err)
	}
	defer rows.Close()

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

func (r *PostgresAdminRepository) ManageUserRequests(conn *pgx.Conn) ([]*models.UnavailableCryptoRequest, error) {
	var requests []*models.UnavailableCryptoRequest

	// Get all unavailable crypto requests
	rows, err := conn.Query(context.Background(), "SELECT crypto_symbol, username, request_message, status, timestamp FROM unavailable_cryptos")
	if err != nil {
		return nil, fmt.Errorf("failed to query unavailable crypto requests: %v", err)
	}
	defer rows.Close()

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

func (r *PostgresAdminRepository) SaveUnavailableCryptoRequest(conn *pgx.Conn, request *models.UnavailableCryptoRequest) error {
	// Insert request
	_, err := conn.Exec(context.Background(), "INSERT INTO unavailable_cryptos (crypto_symbol, username, request_message, status, timestamp) VALUES ($1, $2, $3, $4, $5)",
		request.CryptoSymbol, request.UserName, request.RequestMessage, request.Status, request.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to insert request: %v", err)
	}

	return nil
}
