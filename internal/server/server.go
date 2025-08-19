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
	addr := fmt.Sprintf(":%s", port)

	log.Printf("Starting server at the port: %s\n", port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %s\n", err)
	}
}

/*func Run() {
	port := os.Getenv("HOST_PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    addr,
		Handler: nil,
	}
	go func() {
		log.Printf("Starting server at the port: %s\n", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start server: %s\n", err)
		}
	}()

	sig := <-sigs

	if err := gracefulShutdown(srv, sig); err != nil {
		log.Fatalf("Failed to gracefully shutdown server: %s\n", err)
	}
}*/

/*func gracefulShutdown(srv *http.Server, sig os.Signal) error {
	log.Printf("Received %s signal. Shutting down server...", sig)
	ctx, shutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdown()

	err := srv.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("could not gracefully shutdown server: %w", err)
	}
	return nil
}*/
