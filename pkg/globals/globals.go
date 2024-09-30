package globals

import (
	"context"
	"log"
	"sync"

	"github.com/jackc/pgx/v4"
)

var (
	PgConn     *pgx.Conn
	pgConnOnce sync.Once
	mu         sync.Mutex
)

func GetPgConn() *pgx.Conn {
	mu.Lock()
	defer mu.Unlock()

	pgConnOnce.Do(func() {
		connString := "postgres://postgres:admin_password@localhost:5432/cryptotracker"
		conn, err := pgx.Connect(context.Background(), connString)
		if err != nil {
			log.Fatal("Unable to connect to database:", err)
		}
		PgConn = conn
	})

	return PgConn
}

func ClosePgConn() {
	mu.Lock()
	defer mu.Unlock()

	if PgConn != nil {
		err := PgConn.Close(context.Background())
		if err != nil {
			log.Println("Error closing PostgreSQL connection:", err)
		} else {
			log.Println("PostgreSQL connection closed successfully.")
		}
		PgConn = nil
		pgConnOnce = sync.Once{}
	}
}
