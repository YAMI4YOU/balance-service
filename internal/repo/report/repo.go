package report

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/YAMI4YOU/balance-service/internal/models"
)

type Repo struct {
	conn *pgx.Conn
}

func New(conn *pgx.Conn) *Repo {
	return &Repo{conn: conn}
}

func (db *Repo) FetchReport(ctx context.Context, model models.Report) ([]models.ReportSummary, error) {
	startDate := time.Date(model.Year, time.Month(model.Month), 1, 0, 0, 0, 0, time.UTC)
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
			return nil, fmt.Errorf("failed to scan report record: %w", err)
		}

		summary.Amount = models.NewMoneyFromKopecks(amount)
		summaries = append(summaries, summary)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan report records: %w", err)
	}
	return summaries, nil
}
