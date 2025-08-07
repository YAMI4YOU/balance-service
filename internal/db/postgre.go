package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

type DB struct {
	Conn *pgx.Conn
}

func Init() (*DB, error) {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, os.Getenv("GOOSE_DBSTRING"))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	log.Println("Successfully connected to database")

	return &DB{Conn: conn}, nil
}
