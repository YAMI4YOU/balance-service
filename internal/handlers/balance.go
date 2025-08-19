package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/YAMI4YOU/balance-service/internal/db"
)

type BalanceResponse struct {
	UserID  int   `json:"user_id"`
	Balance int64 `json:"balance"`
}

type BHandler struct {
	store *db.DB
}

func BalanceH(store *db.DB) *BHandler {
	return &BHandler{store: store}
}

func (h *BHandler) BalanceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Printf("Got request for balance: %s", r.URL)

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user_id parameter", http.StatusBadRequest)
		return
	}

	balance, err := h.store.GetBalance(r.Context(), userID)
	if err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, BalanceResponse{UserID: userID, Balance: balance}, http.StatusOK)
}

// writeJSON отправляет JSON-ответ
func writeJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	d, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "ошибка сериализации JSON", http.StatusInternalServerError)
	}

	_, err = w.Write(d)
	if err != nil {
		log.Printf("error writing response: %w", err)
		return
	}
}
