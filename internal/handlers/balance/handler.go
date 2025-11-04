package balance

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/YAMI4YOU/balance-service/internal/handlers"
	"github.com/YAMI4YOU/balance-service/internal/models"
)

type response struct {
	UserID  int     `json:"user_id"`
	Balance float64 `json:"balance"`
}

type Handler struct {
	Store repo
}

func NewHandler(store repo) *Handler {
	return &Handler{Store: store}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Printf(`Got request "GET" for balance: %s`, r.URL)

	userIDFromURL := r.URL.Query().Get("user_id")
	if userIDFromURL == "" {
		http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDFromURL)
	if err != nil {
		http.Error(w, "Invalid user_id parameter", http.StatusBadRequest)
		return
	}

	balance, err := h.Store.FetchBalance(r.Context(), userID)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handlers.WriteJSON(w,
		response{
			UserID:  userID,
			Balance: balance.Rubles()},
		http.StatusOK)
}
