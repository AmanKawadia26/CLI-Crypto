package repositories_test

import (
	"context"
	"cryptotracker/internal/repositories"
	//"cryptotracker/models"
	//"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserProfile(t *testing.T) {
	conn, cleanup, err := setupTestDB()
	require.NoError(t, err)
	defer cleanup()

	repo := repositories.NewPostgresUserRepository(conn)

	_, err = conn.Exec(context.Background(), "INSERT INTO users (userid, username, email, mobile, isadmin, role) VALUES ($1, $2, $3, $4, $5, $6)",
		"0", "testuser1", "testuser@example.com", 1234567890, false, "user")
	require.NoError(t, err)

	user, err := repo.GetUserProfile("testuser1")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	//assert.Equal(t, "0", user.UserID)
	assert.Equal(t, "testuser1", user.Username)
	assert.Equal(t, "testuser@example.com", user.Email)
	assert.Equal(t, 1234567890, user.Mobile)
	assert.Equal(t, "user", user.Role)

	user, err = repo.GetUserProfile("nonexistentuser")
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "failed to get user profile: no rows in result set", err.Error())

	_, err = conn.Exec(context.Background(), "DROP TABLE users")
	require.NoError(t, err)

	user, err = repo.GetUserProfile("testuserabc")
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "failed to get user profile")
}
