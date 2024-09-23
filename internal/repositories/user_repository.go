package repositories

import (
	"context"
	"cryptotracker/models"
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
	query := `
		SELECT username, email, mobile, role
		FROM users
		WHERE username = $1
	`

	var user models.User
	err := repo.conn.QueryRow(context.Background(), query, username).Scan(
		&user.Username,
		&user.Email,
		&user.Mobile,
		&user.Role,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %v", err)
	}
	return &user, nil
}
