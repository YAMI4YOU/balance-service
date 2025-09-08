package reservation

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
	UserID    int     `json:"user_id"`
	ServiceID int     `json:"service_id"`
	OrderID   int     `json:"order_id"`
	Amount    float64 `json:"amount"`
}

func (req *Request) Validate() error {
	if req.UserID <= 0 {
		return errors.New("invalid user_id")
	}
	if req.ServiceID <= 0 {
		return errors.New("invalid service_id")
	}
	if req.OrderID <= 0 {
		return errors.New("invalid order_id")
	}
	if req.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	return nil
}

type repo interface {
	ReserveFunds(ctx context.Context, req models.Reservation) error
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

	err := req.Validate()
	if err != nil {
		log.Printf("Reservation request validation failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.store.ReserveFunds(r.Context(), models.Reservation{
		UserID:    req.UserID,
		ServiceID: req.ServiceID,
		OrderID:   req.OrderID,
		Amount:    models.NewMoneyFromRubles(req.Amount),
	}); err != nil {
		switch {
		case errors.Is(err, models.ErrUserNotFound):
			log.Printf("User %d not found", req.UserID)
			http.Error(w, "User not found", http.StatusNotFound)
		case errors.Is(err, models.ErrInsufficientFunds):
			log.Printf("Insufficient funds for user %d", req.UserID)
			http.Error(w, "Insufficient funds", http.StatusPaymentRequired)
		case errors.Is(err, models.ErrConflictReservation):
			log.Printf("Conflict reservation for user %d", req.UserID)
			http.Error(w, "Reservation already exists for other user", http.StatusConflict)
		default:
			log.Printf("Reservation failed: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("Successfully reserved funds for user: %d, amount: %.2f",
		req.UserID, req.Amount)

	handlers.WriteJSON(w, map[string]string{
		"status":  "success",
		"message": "Funds reserved successfully"},
		http.StatusCreated)
}
