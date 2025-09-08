package reportservice

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/YAMI4YOU/balance-service/internal/models"
)

type recordReport struct {
	UserID    int64  `csv:"user_id"`
	ServiceID int64  `csv:"service_id"`
	OrderID   int64  `csv:"order_id"`
	AmountRub string `csv:"amount_rub"`
	Operation string `csv:"operation"`
	CreatedAt string `csv:"created_at"`
}

type Repository interface {
	MonthlyReport(ctx context.Context, req models.MonthlyReportRequest) ([]models.MonthlyReportRecord, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GenerateMonthlyReport(ctx context.Context, req models.MonthlyReportRequest) (string, error) {
	records, err := s.repo.MonthlyReport(ctx, req)
	if err != nil {
		return "", fmt.Errorf("generate monthly reportservice error: %w", err)
	}

	filename := fmt.Sprintf("report_%d_%02d_%s.csv", req.Year, req.Month,
		time.Now().Format("20060102_150405"))
	filepath := filepath.Join("reports", filename)

	if err := os.MkdirAll("reports", 0777); err != nil {
		return "", fmt.Errorf("create reportservice directory error: %w", err)
	}

	file, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("create reportservice file error: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"user_id", "service_id", "order_id", "amount_rub", "operation", "created_at"}
	if err := writer.Write(header); err != nil {
		return "", fmt.Errorf("write header to reportservice file error: %w", err)
	}

	for _, record := range records {

		reportRecord := convertToRecordReport(record)

		row := []string{
			strconv.FormatInt(reportRecord.UserID, 10),
			strconv.FormatInt(reportRecord.ServiceID, 10),
			strconv.FormatInt(reportRecord.OrderID, 10),
			reportRecord.AmountRub,
			reportRecord.Operation,
			reportRecord.CreatedAt,
		}

		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("write row to reportservice file error: %w", err)
		}
	}

	return filename, nil
}

func convertToRecordReport(record models.MonthlyReportRecord) recordReport {
	return recordReport{
		UserID:    record.UserID,
		ServiceID: record.ServiceID,
		OrderID:   record.OrderID,
		AmountRub: fmt.Sprintf("%.2f", record.Amount.Rubles()),
		Operation: record.Operation,
		CreatedAt: record.CreatedAt,
	}
}
