package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

var (
	conn *pgx.Conn
)

// Init initializes the database connection
func Init() error {
	var err error

	// get the database url from environment variable
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return fmt.Errorf("database url not found")
	}

	// create a new database connection
	conn, err = pgx.Connect(context.Background(), url)
	if err != nil {
		return err
	}

	return nil
}

// Close closes the database connection
func Close() {
	if conn != nil {
		conn.Close(context.Background())
	}
}
