package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func Run() {
	port := os.Getenv("HOST_PORT")
	if port == "" {
		port = "8080"
	}

	// Запуск сервера
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting server at the port: %s\n", port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %s\n", err)
	}
}
