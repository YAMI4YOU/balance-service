package report

import (
	"context"

	"github.com/YAMI4YOU/balance-service/internal/models"
)

type service interface {
	GenerateReport(ctx context.Context, model models.Report) (string, error)
}
