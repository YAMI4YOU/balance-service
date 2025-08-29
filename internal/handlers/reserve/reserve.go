package reserve

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/YAMI4YOU/balance-service/internal/handlers"
	"github.com/YAMI4YOU/balance-service/internal/models"
)

type Request struct {
	UserID    int   `json:"user_id"`
	ServiceID int   `json:"service_id"`
	OrderID   int   `json:"order_id"`
	Amount    int64 `json:"amount"`
}

type repo interface {
	ReserveFunds(ctx context.Context, req models.Reserve) error
}

type Handler struct {
	store repo
}

func NewHandler(store repo) *Handler {
	return &Handler{store: store}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Printf(`Got reguest "POST" for reservation: %s`, r.URL)

	var req Request
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

	if err := h.store.ReserveFunds(r.Context(), models.Reserve{
		UserID:    req.UserID,
		ServiceID: req.ServiceID,
		OrderID:   req.OrderID,
		Amount:    req.Amount,
	}); err != nil {
		switch {
		case errors.Is(err, models.ErrUserNotFound):
			http.Error(w, "User not found", http.StatusNotFound)
		case errors.Is(err, models.ErrInsufficientFunds):
			http.Error(w, "Insufficient funds", http.StatusPaymentRequired)
		default:
			log.Printf("Reserve failed: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	handlers.WriteJSON(w, map[string]string{
		"status":  "success",
		"message": "Funds reserved successfully"},
		http.StatusCreated)
}
