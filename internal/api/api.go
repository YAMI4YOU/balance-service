package api

import (
	"net/http"

	"balance-service/internal/db"
	"balance-service/internal/handlers"
)

func Init(db *db.DB) {
	balanceHandler := handlers.BalanceH(db)
	depositHandler := handlers.DepositH(db)

	http.HandleFunc("/balance", balanceHandler.BalanceHandler)
	http.HandleFunc("/deposit", depositHandler.DepositHandler)
}
