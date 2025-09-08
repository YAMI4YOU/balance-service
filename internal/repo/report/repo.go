package report

import (
	"context"
	"fmt"
	"time"

	"github.com/YAMI4YOU/balance-service/internal/models"
	"github.com/jackc/pgx/v5"
)

type Repo struct {
	conn *pgx.Conn
}

func New(conn *pgx.Conn) *Repo {
	return &Repo{conn: conn}
}

func (db *Repo) MonthlyReport(ctx context.Context, req models.MonthlyReportRequest) ([]models.MonthlyReportRecord, error) {
	if req.Year < 2000 || req.Year > 2100 {
		return nil, fmt.Errorf("year out of range: %d", req.Year)
	}

	if req.Month < 1 || req.Month > 12 {
		return nil, fmt.Errorf("month out of range: %d", req.Month)
	}

	startDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	query := `
		SELECT user_id, service_id, order_id, amount, operation_type, created_at
		FROM accounting_report
		WHERE created_at >= $1 AND created_at < $2
		ORDER BY created_at ASC
	`

	rows, err := db.conn.Query(ctx, query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("select accounting reportservice failed: %w", err)
	}

	defer rows.Close()

	var records []models.MonthlyReportRecord
	for rows.Next() {
		var record models.MonthlyReportRecord
		var amount int64
		var createdAt time.Time

		err := rows.Scan(
			&record.UserID,
			&record.ServiceID,
			&record.OrderID,
			&amount,
			&record.Operation,
			&createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reportservice record: %w", err)
		}

		record.Amount = models.NewMoneyFromKopecks(amount)
		record.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan reportservice records: %w", err)
	}
	return records, nil
}
