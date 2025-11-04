package file

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/YAMI4YOU/balance-service/internal/models"
)

type Writer struct {
}

func NewWriter() *Writer {
	return &Writer{}
}

func (f *Writer) Save(model models.Report, summaries []models.ReportSummary) (string, error) {
	filename := fmt.Sprintf("report_%d_%02d_%s.csv", model.Year, model.Month, time.Now().Format("20060102_150405"))
	pathOfFile := filepath.Join("reports", filename)

	if err := os.MkdirAll("reports", 0777); err != nil {
		return "", fmt.Errorf("create report directory error: %w", err)
	}

	file, err := os.Create(pathOfFile)
	if err != nil {
		return "", fmt.Errorf("create report service file error: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	header := []string{"service_id", "total_amount_rub"}
	if err := writer.Write(header); err != nil {
		return "", fmt.Errorf("write header error: %w", err)
	}

	for _, summary := range summaries {
		row := []string{
			strconv.FormatInt(summary.ServiceID, 10),
			fmt.Sprintf("%.2f", summary.Amount.Rubles()),
		}
		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("write row error: %w", err)
		}
	}

	return filename, err
}
