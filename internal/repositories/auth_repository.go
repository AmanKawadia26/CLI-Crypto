package repositories

import (
	"context"
	"cryptotracker/models"
	"cryptotracker/pkg/config"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"log"
)

type AuthRepository interface {
	LoginDBRepository(username string) (*models.User, error)
	SignupDBRepository(user *models.User) error
}

type PostgresAuthRepository struct {
	conn *pgx.Conn
}

func NewPostgresAuthRepository(conn *pgx.Conn) *PostgresAuthRepository {
	return &PostgresAuthRepository{conn: conn}
}

// LoginDBRepository handles the database interaction for user login.
func (r *PostgresAuthRepository) LoginDBRepository(username string) (*models.User, error) {

	var user models.User

	// Define the columns and conditions for the query
	columns := []string{"username", "password", "email", "mobile", "role", "isadmin"}
	condition := "username = $1"

	// Build the query to fetch user information
	query, err := config.BuildSelectQuery(columns, "users", condition)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	// Execute the query to fetch the user information
	err = r.conn.QueryRow(context.Background(), query, username).
		Scan(&user.Username, &user.Password, &user.Email, &user.Mobile, &user.Role, &user.IsAdmin)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err // Return other DB-related errors
	}

	return &user, nil
}

func (r *PostgresAuthRepository) SignupDBRepository(user *models.User) error {
	// Define columns for the SELECT query
	selectColumns := []string{"username"}
	selectCondition := "username = $1"

	// Use BuildSelectQuery to create a SELECT query for checking existing user
	selectQuery, err := config.BuildSelectQuery(selectColumns, "users", selectCondition)
	if err != nil {
		return fmt.Errorf("failed to build select query: %v", err)
	}

	// Check if the user already exists
	var existingUser models.User
	err = r.conn.QueryRow(context.Background(), selectQuery, user.Username).Scan(&existingUser.Username)
	if err == nil {
		return errors.New("user already exists")
	} else if err != pgx.ErrNoRows {
		log.Fatalf("failed to query user: %v", err)
		return err
	}

	// Define columns for the INSERT query
	insertColumns := []string{"username", "email", "mobile", "password", "role", "isadmin"}

	// Use BuildInsertQuery to create an INSERT query for user registration
	insertQuery, err := config.BuildInsertQuery("users", insertColumns)
	if err != nil {
		return fmt.Errorf("failed to build insert query: %v", err)
	}

	// Insert the new user
	_, err = r.conn.Exec(context.Background(), insertQuery, user.Username, user.Email, user.Mobile, user.Password, user.Role, user.IsAdmin)
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
		return err
	}

	return nil
}
