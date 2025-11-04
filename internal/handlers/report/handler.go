package report

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/YAMI4YOU/balance-service/internal/handlers"
	"github.com/YAMI4YOU/balance-service/internal/models"
)

type Handler struct {
	service service
}

func NewHandler(s service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Printf(`Got reguest "GET" for recognize revenue: %s`, r.URL)

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		log.Printf("Report request validation failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filepath, err := h.service.GenerateReport(r.Context(), models.Report{
		Year:  req.Year,
		Month: req.Month,
	})

	if err != nil {
		log.Printf("failed to generate report: %v", err)
		http.Error(w, "failed to generate report", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"status":  "success",
		"file":    fmt.Sprintf("/reports/%s", filepath),
		"message": fmt.Sprintf("Report for %d-%02d generated successfully", req.Year, req.Month),
	}

	log.Printf("Successfully created report for year: %d and month: %d",
		req.Year, req.Month)

	handlers.WriteJSON(w, response, http.StatusOK)
}
