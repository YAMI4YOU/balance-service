package main

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/YAMI4YOU/balance-service/config"
	"github.com/YAMI4YOU/balance-service/internal/configure"
	"github.com/YAMI4YOU/balance-service/internal/handlers/balance"
	"github.com/YAMI4YOU/balance-service/internal/handlers/deposit"
	"github.com/YAMI4YOU/balance-service/internal/handlers/reserve"
	balanceRepo "github.com/YAMI4YOU/balance-service/internal/repo/balance"
	depositRepo "github.com/YAMI4YOU/balance-service/internal/repo/deposit"
	reserveRepo "github.com/YAMI4YOU/balance-service/internal/repo/reservation"
	"github.com/YAMI4YOU/balance-service/internal/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Application initialization failed: %v", r)
			log.Printf("Stack trace: %s", debug.Stack())
		}
	}()

	connection := configure.MustInitDB(ctx, cfg.DBUrl)
	defer func() {
		ctxClose, cancelClose := context.WithTimeout(context.Background(), time.Second)
		defer cancelClose()

		_ = connection.Close(ctxClose)
	}()

	// <-- Repo
	repoBalance := balanceRepo.New(connection)
	repoDeposit := depositRepo.New(connection)
	repoReserve := reserveRepo.New(connection)
	// Repo -->

	// <-- Handle
	balanceHandler := balance.NewHandler(repoBalance)
	depositHandler := deposit.NewHandler(repoDeposit)
	reserveHandler := reserve.NewHandler(repoReserve)
	//  Handle -->

	http.HandleFunc("/balance", balanceHandler.Handle)
	http.HandleFunc("/deposit", depositHandler.Handle)
	http.HandleFunc("/reserve", reserveHandler.Handle)

	server.Run(ctx, cancel, cfg.HostPort)
}
