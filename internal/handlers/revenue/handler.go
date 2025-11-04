package revenue

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/YAMI4YOU/balance-service/internal/handlers"
	"github.com/YAMI4YOU/balance-service/internal/handlers/request"
	"github.com/YAMI4YOU/balance-service/internal/models"
)

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
	log.Printf(`Got reguest "POST" for recognize revenue: %s`, r.URL)

	var req request.ReservationAndRevenueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		log.Printf("Revenue request validation failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.store.RecognizeRevenue(r.Context(), models.Revenue{
		UserID:    req.UserID,
		ServiceID: req.ServiceID,
		OrderID:   req.OrderID,
		Amount:    models.NewMoneyFromRubles(req.Amount),
	}); err != nil {
		switch {
		case errors.Is(err, models.ErrReservationNotFound):
			log.Printf("Reservation not found for user: %d, service: %d, order: %d",
				req.UserID, req.ServiceID, req.OrderID)
			http.Error(w, "Reservation not found", http.StatusNotFound)
		case errors.Is(err, models.ErrReservationNotReserved):
			log.Printf("Reservation is not in reserved status for user: %d", req.UserID)
			http.Error(w, "Reservation is not in reserved status", http.StatusBadRequest)
		case errors.Is(err, models.ErrUserNotFound):
			log.Printf("User not found: %d", req.UserID)
			http.Error(w, "User not found", http.StatusNotFound)
		case errors.Is(err, models.ErrUserMismatch):
			log.Printf("User mismatch: %d", req.UserID)
			http.Error(w, "User mismatch", http.StatusBadRequest)
		case errors.Is(err, models.ErrAmountMismatch):
			log.Printf("Amount mismatch: %v", req.Amount)
			http.Error(w, "Amount mismatch", http.StatusBadRequest)
		default:
			log.Printf("Revenue recognition failed: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("Successfully recognized revenue for user: %d, amount: %.2f",
		req.UserID, req.Amount)

	handlers.WriteJSON(w, map[string]string{
		"status":  "success",
		"message": "Funds reserved successfully"},
		http.StatusCreated)
}
