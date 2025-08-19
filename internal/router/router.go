package router

import (
	"net/http"

	"github.com/YAMI4YOU/balance-service/internal/db"
	"github.com/YAMI4YOU/balance-service/internal/handlers"
)

func Init(db *db.DB) {
	balanceHandler := handlers.BalanceH(db)
	depositHandler := handlers.NewDepositH(db)
	reserveHandler := handlers.ReserveH(db)

	http.HandleFunc("/balance", balanceHandler.BalanceHandler)
	http.HandleFunc("/deposit", depositHandler.DepositHandler)
	http.HandleFunc("/reserve", reserveHandler.ReserveHandler)
}
