package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

type DB struct {
	Conn *pgx.Conn
}

// сделать возврат Conn
func MustInit() *DB {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, os.Getenv("GOOSE_DBSTRING"))
	if err != nil {
		panic("unable to connect to database")
	}

	// сделать ретрай 3 раза пинг
	if err := conn.Ping(ctx); err != nil {
		panic("unable to ping database")
	}

	log.Println("Successfully connected to database")

	var version string
	err = conn.QueryRow(context.Background(), "SELECT version()").Scan(&version)
	if err != nil {
		log.Printf("Version query failed: %v", err)
	}
	log.Println("PostgreSQL version:", version)

	return &DB{Conn: conn}
}

func (db *DB) CloseDB() {
	if db.Conn == nil {
		return
	}
	if err := db.Conn.Close(context.Background()); err != nil {
		log.Printf("Failed to close connection: %v\n", err)
	} else {
		log.Println("Successfully closed connection")
	}
	db.Conn = nil
}
