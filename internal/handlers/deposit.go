package handlers

import (
	"balance-service/internal/db"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type DHandler struct {
	store *db.DB
}

func DepositH(store *db.DB) *DHandler {
	return &DHandler{store: store}
}

func (h *DHandler) DepositHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	log.Printf("Headers: %+v\n", r.Header)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %s", err)
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}
	fmt.Printf("body: %s\n", string(body))
	defer r.Body.Close()

	var req db.DepositRequest
	if err = json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.UserID <= 0 {
		http.Error(w, "Invalid user_id parameter", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "Invalid amount parameter", http.StatusBadRequest)
		return
	}

	if err = h.store.Deposit(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed data: %+v\n", req)

	writeJSON(w, map[string]string{"status": "success"}, http.StatusCreated)
}
