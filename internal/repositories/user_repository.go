package repositories

import (
	"context"
	"cryptotracker/models"
	"cryptotracker/pkg/config"
	"fmt"
	"github.com/jackc/pgx/v4"
	//"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	GetUserProfile(username string) (*models.User, error)
}

type PostgresUserRepository struct {
	conn *pgx.Conn
}

func NewPostgresUserRepository(conn *pgx.Conn) UserRepository {
	return &PostgresUserRepository{conn: conn}
}

func (repo *PostgresUserRepository) GetUserProfile(username string) (*models.User, error) {
	// Define the columns to select
	columns := []string{"username", "email", "mobile", "role"}
	// Define the condition
	condition := "username = $1"

	// Build the SELECT query using BuildSelectQuery
	query, err := config.BuildSelectQuery(columns, "users", condition)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	var user models.User
	err = repo.conn.QueryRow(context.Background(), query, username).Scan(
		&user.Username,
		&user.Email,
		&user.Mobile,
		&user.Role,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("no rows in result set")
		}
		return nil, fmt.Errorf("failed to get user profile: %v", err)
	}
	return &user, nil
}
