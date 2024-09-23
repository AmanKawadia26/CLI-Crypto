package globals_test

import (
	"context"
	"testing"

	"cryptotracker/pkg/globals" // Update with the actual module path
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"
)

const testConnString = "postgres://postgres:admin_password@localhost:5432/cryptotracker_test" // Use a separate test database

func setupTestDB(t *testing.T) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), testConnString)
	require.NoError(t, err, "failed to connect to test database")

	_, err = conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS test_table (
			id SERIAL PRIMARY KEY,
			data TEXT
		)
	`)
	require.NoError(t, err, "failed to create test table")

	return conn
}

func teardownTestDB(t *testing.T, conn *pgx.Conn) {
	_, err := conn.Exec(context.Background(), `DROP TABLE IF EXISTS test_table`)
	require.NoError(t, err, "failed to drop test table")
	conn.Close(context.Background())
}

func TestGetPgConn(t *testing.T) {
	conn := setupTestDB(t)
	defer teardownTestDB(t, conn)

	// Ensure GetPgConn returns a non-nil connection
	conn1 := globals.GetPgConn()
	require.NotNil(t, conn1, "expected a non-nil connection")

	// Optionally test that the connection is indeed valid
	err := conn1.Ping(context.Background())
	require.NoError(t, err, "failed to ping the database")

	// Check that subsequent calls to GetPgConn return the same connection
	conn2 := globals.GetPgConn()
	require.Equal(t, conn1, conn2, "expected the same connection instance")

	// Close the connection after the test to avoid resource leaks
	defer globals.ClosePgConn()
}

func TestClosePgConn(t *testing.T) {
	conn := globals.GetPgConn()
	require.NotNil(t, conn, "expected a non-nil connection before closing")

	// Close the connection
	globals.ClosePgConn()

	// Ensure the connection was closed by attempting to get a new one
	conn1 := globals.GetPgConn()
	require.NotNil(t, conn1, "expected a non-nil connection after closing and reopening")

	// Test that the new connection is valid
	err := conn1.Ping(context.Background())
	require.NoError(t, err, "failed to ping the database after reopening")

	// Close the connection again to clean up
	defer globals.ClosePgConn()
}

func TestClosePgConnNoConnection(t *testing.T) {
	// Close connection when none is open
	globals.ClosePgConn()

	// Ensure no error is logged
	// This test mainly ensures the ClosePgConn function handles the nil case gracefully
}
