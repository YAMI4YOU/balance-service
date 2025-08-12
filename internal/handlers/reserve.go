package handlers

import (
	"balance-service/internal/db"
	"encoding/json"
	"log"
	"net/http"
)

type RHandler struct {
	store *db.DB
}

func ReserveH(store *db.DB) *RHandler {
	return &RHandler{store: store}
}

func (h *RHandler) ReserveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req db.Reserve
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.UserID <= 0 {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}
	if req.ServiceID <= 0 {
		http.Error(w, "Invalid service_id", http.StatusBadRequest)
		return
	}
	if req.OrderID <= 0 {
		http.Error(w, "Invalid order_id", http.StatusBadRequest)
		return
	}
	if req.Amount <= 0 {
		http.Error(w, "Amount must be greater than 0", http.StatusBadRequest)
		return
	}

	if err := h.store.ReserveFunds(r.Context(), req); err != nil {
		switch err.Error() {
		case "user not found":
			http.Error(w, "User not found", http.StatusNotFound)
		case "insufficient funds":
			http.Error(w, "Insufficient funds", http.StatusPaymentRequired)
		default:
			log.Printf("Reserve failed: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, map[string]string{
		"status":  "success",
		"message": "Funds reserved successfully"},
		http.StatusCreated)
}
