package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	server *http.Server
}

func NewServer(hostPort string) *Server {
	addr := ":" + hostPort
	return &Server{
		server: &http.Server{
			Addr: addr,
		},
	}
}

func (s *Server) Start() {

	fmt.Printf("Starting server at %s\n", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Error starting server: %s\n", err)
	}

}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) Addr() string {
	return s.server.Addr
}
