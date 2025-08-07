package main

import (
	"balance-service/internal/api"
	"context"
	"log"

	"balance-service/internal/db"
	"balance-service/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: No .env file found, using system environment variables")
	}

	conn, err := db.Init()
	if err != nil {
		log.Fatalf("Failed to connect to database: %s\n", err)
	}
	defer func() {
		if conn.Conn != nil {
			conn.Conn.Close(context.Background())
		}
	}()

	var version string
	err = conn.Conn.QueryRow(context.Background(), "SELECT version()").Scan(&version)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	log.Println("PostgreSQL version:", version)

	api.Init(conn)

	server.Run()
}
