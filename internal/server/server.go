package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(ctx context.Context, cancel context.CancelFunc, hostPort string) {
	addr := ":" + hostPort
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr: addr,
	}
	go func() {
		fmt.Println("Listening on " + addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start server: %s\n", err)
		}
	}()

	sig := <-sigs
	cancel()

	if err := gracefulShutdown(srv, sig); err != nil {
		log.Fatalf("Failed to gracefully shutdown server: %s\n", err)
	}
}

func gracefulShutdown(srv *http.Server, sig os.Signal) error {
	log.Printf("Received %s signal. Shutting down server...", sig)
	ctx, shutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdown()

	err := srv.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("could not gracefully shutdown server: %w", err)
	}
	return nil
}
