package deposit

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/YAMI4YOU/balance-service/internal/handlers"
	"github.com/YAMI4YOU/balance-service/internal/models"
)

type request struct {
	UserID  int     `json:"user_id"`
	Balance float64 `json:"balance"`
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
	}
	log.Printf(`Got reguest "POST" for deposit: %s`, r.URL)
	log.Printf("Headers: %+v\n", r.Header)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %s", err)
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}
	fmt.Printf("body: %s\n", string(body))
	defer r.Body.Close()

	var req request
	if err = json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.UserID <= 0 {
		http.Error(w, "Invalid user_id parameter", http.StatusBadRequest)
		return
	}

	if req.Balance <= 0 {
		http.Error(w, "Invalid amount parameter", http.StatusBadRequest)
		return
	}

	if err = h.store.MakeDeposit(r.Context(), models.Deposit{
		UserID:  req.UserID,
		Balance: models.NewMoneyFromRubles(req.Balance),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed data: %+v\n", req)

	handlers.WriteJSON(w, map[string]string{"status": "success"}, http.StatusCreated)
}
