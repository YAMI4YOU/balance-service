package revenue

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/YAMI4YOU/balance-service/internal/models"
	"github.com/jackc/pgx/v5"
)

type Repo struct {
	conn *pgx.Conn
}

func New(conn *pgx.Conn) *Repo {
	return &Repo{conn: conn}
}

func (db *Repo) RecognizeRevenue(ctx context.Context, revenue models.Revenue) error {
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("recognize revenue transaction begin: %w", err)
	}
	defer tx.Rollback(ctx)

	// Может это сделать структурой?
	var reservationID int64
	var currentStatus string
	var reservedAmount int64
	var reservationUserID int

	err = tx.QueryRow(ctx, `
		SELECT id, user_id, amount, status
		FROM reservation
		WHERE service_id = $1 and order_id = $2
		FOR UPDATE
	`, revenue.ServiceID,
		revenue.OrderID).Scan(&reservationID, &reservationUserID, &reservedAmount, &currentStatus)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ErrReservationNotFound
		}
		return fmt.Errorf("recognize revenue transaction row: %w", err)
	}

	if reservationUserID != revenue.UserID {
		log.Printf("reservation belongs to user %d, not to %d: %v",
			reservationUserID, revenue.UserID, models.ErrUserMismatch)
		return models.ErrUserMismatch
	}

	if currentStatus != string(models.ReservationStatusReserved) {
		return models.ErrReservationNotReserved
	}

	if revenue.Amount.Kopecks() != reservedAmount {
		log.Printf("recognize revenue requsted %d does not match reserved amount %d",
			revenue.Amount, reservedAmount)
		return models.ErrAmountMismatch
	}

	_, err = tx.Exec(ctx, `
		UPDATE reservation
		SET status = $1
		WHERE id = $2
	`, models.ReservationStatusConfirmed, reservationID)

	if err != nil {
		return fmt.Errorf("recognize revenue transaction update: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO accounting_report (user_id, service_id, order_id, amount, operation_type, created_at)
		    VALUES ($1, $2, $3, $4, $5, NOW())
	`, revenue.UserID, revenue.ServiceID, revenue.OrderID, revenue.Amount.Kopecks(), "revenue_recognition")

	if err != nil {
		return fmt.Errorf("create accounting report failed: %w", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("recognize revenue transaction commit: %w", err)
	}

	return nil
}
