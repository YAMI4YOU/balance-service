//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package report

import (
	"context"

	"github.com/YAMI4YOU/balance-service/internal/models"
)

type repository interface {
	FetchReport(ctx context.Context, req models.Report) ([]models.ReportSummary, error)
}

type file interface {
	Save(model models.Report, summaries []models.ReportSummary) (string, error)
}
