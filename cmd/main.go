package main

import (
	"log"
	"runtime/debug"

	"github.com/YAMI4YOU/balance-service/internal/db"
	"github.com/YAMI4YOU/balance-service/internal/router"
	"github.com/YAMI4YOU/balance-service/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: No .env file found, using system environment variables")
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Application initialization failed: %v", r)
			log.Printf("Stack trace: %s", debug.Stack())
		}
	}()

	connection := db.MustInit()
	defer connection.CloseDB()

	router.Init(connection)

	server.Run()
}
