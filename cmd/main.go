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
	"github.com/YAMI4YOU/balance-service/internal/handlers/report"
	"github.com/YAMI4YOU/balance-service/internal/handlers/reservation"
	"github.com/YAMI4YOU/balance-service/internal/handlers/revenue"
	balanceRepo "github.com/YAMI4YOU/balance-service/internal/repo/balance"
	depositRepo "github.com/YAMI4YOU/balance-service/internal/repo/deposit"
	reportRepo "github.com/YAMI4YOU/balance-service/internal/repo/report"
	reserveRepo "github.com/YAMI4YOU/balance-service/internal/repo/reservation"
	revenueRepo "github.com/YAMI4YOU/balance-service/internal/repo/revenue"
	"github.com/YAMI4YOU/balance-service/internal/server"
	"github.com/YAMI4YOU/balance-service/internal/service/reportservice"
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
	repoRevenue := revenueRepo.New(connection)
	repoReport := reportRepo.New(connection)
	// Repo -->

	// <-- Service
	reportService := reportservice.NewService(repoReport)
	// Service -->

	// <-- Handle
	balanceHandler := balance.NewHandler(repoBalance)
	depositHandler := deposit.NewHandler(repoDeposit)
	reserveHandler := reservation.NewHandler(repoReserve)
	revenueHandler := revenue.NewHandler(repoRevenue)
	reportHandler := report.NewHandler(reportService)
	//  Handle -->

	// <-- Route
	http.HandleFunc("/balance", balanceHandler.Handle)
	http.HandleFunc("/deposit", depositHandler.Handle)
	http.HandleFunc("/reserve", reserveHandler.Handle)
	http.HandleFunc("/revenue", revenueHandler.Handle)
	http.HandleFunc("/report", reportHandler.Handle)
	// Route -->

	server.Run(ctx, cancel, cfg.HostPort)
}
