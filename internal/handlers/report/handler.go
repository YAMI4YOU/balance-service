package report

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/YAMI4YOU/balance-service/internal/handlers"
	"github.com/YAMI4YOU/balance-service/internal/models"
	"github.com/YAMI4YOU/balance-service/internal/service/reportservice"
)

type request struct {
	Year  int `json:"year"`
	Month int `json:"month"`
}

func (req *request) Validate() error {
	if req.Year < 2000 || req.Year > 2100 {
		return errors.New("year out of range")
	}
	if req.Month < 1 || req.Month > 12 {
		return errors.New("month out of range")
	}

	return nil
}

type Handler struct {
	service *reportservice.Service
}

func NewHandler(s *reportservice.Service) *Handler {
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

	filepath, err := h.service.GenerateMonthlyReport(r.Context(), models.MonthlyReportRequest{
		Year:  req.Year,
		Month: req.Month,
	})

	if err != nil {
		log.Printf("failed to generate report service: %v", err)
		http.Error(w, "failed to generate report service", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"status":  "success",
		"file":    fmt.Sprintf("/reports/%s", filepath),
		"message": fmt.Sprintf("Report for %d-%02d generated successfully", req.Year, req.Month),
	}

	log.Printf("Successfully created reportservice for year: %d and month: %d",
		req.Year, req.Month)

	handlers.WriteJSON(w, response, http.StatusOK)
}
