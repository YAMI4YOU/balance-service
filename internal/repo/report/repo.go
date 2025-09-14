package report

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/YAMI4YOU/balance-service/internal/models"
)

const (
	startYear = 2000
	endYear   = 2100

	firstMonth = 1
	lastMonth  = 12
)

type Repo struct {
	conn *pgx.Conn
}

func New(conn *pgx.Conn) *Repo {
	return &Repo{conn: conn}
}

func (db *Repo) MonthlyReport(ctx context.Context, req models.MonthlyReportRequest) ([]models.ReportSummary, error) {
	if req.Year < startYear || req.Year > endYear {
		return nil, fmt.Errorf("year out of range: %d", req.Year)
	}

	if req.Month < firstMonth || req.Month > lastMonth {
		return nil, fmt.Errorf("month out of range: %d", req.Month)
	}

	startDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	query := `
		SELECT service_id, SUM(amount) as total_amount
		FROM accounting_report
		WHERE created_at >= $1 AND created_at < $2
		GROUP BY service_id
		ORDER BY service_id
	`

	rows, err := db.conn.Query(ctx, query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("select accounting report service failed: %w", err)
	}

	defer rows.Close()

	var summaries []models.ReportSummary
	for rows.Next() {
		var summary models.ReportSummary
		var amount int64

		err := rows.Scan(
			&summary.ServiceID,
			&amount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reportservice record: %w", err)
		}

		summary.Amount = models.NewMoneyFromKopecks(amount)
		summary.CreatedAt = fmt.Sprintf("%d-%02d", req.Year, req.Month)
		summaries = append(summaries, summary)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan reportservice records: %w", err)
	}
	return summaries, nil
}
