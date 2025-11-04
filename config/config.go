package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HostPort string `env:"PORT"`
	DBUrl    string `env:"GOOSE_DBSTRING"`
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: No .env file found, using system environment variables")
	}

	cfg := &Config{
		HostPort: os.Getenv("PORT"),
		DBUrl:    os.Getenv("GOOSE_DBSTRING"),
	}

	flag.StringVar(&cfg.HostPort, "port", cfg.HostPort, "The host:port to listen on")
	flag.StringVar(&cfg.DBUrl, "dburl", cfg.DBUrl, "The database url")
	flag.Parse()

	if cfg.DBUrl == "" {
		return nil, fmt.Errorf("databace URL is required")
	}
	return cfg, nil
}
