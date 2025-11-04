package reservation

import (
	"context"
	"errors"
	"fmt"

	"github.com/YAMI4YOU/balance-service/internal/models"
	"github.com/jackc/pgx/v5"
)

type Repo struct {
	conn *pgx.Conn
}

func New(conn *pgx.Conn) *Repo {
	return &Repo{conn: conn}
}

func (db *Repo) ReserveFunds(ctx context.Context, model models.Reservation) error {
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("couldn't start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var currentBalance int64
	err = tx.QueryRow(ctx, `
		SELECT balance 
		FROM wallet 
		WHERE user_id = $1 
		FOR UPDATE
	`, model.UserID).Scan(&currentBalance)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ErrUserNotFound
		}
		return fmt.Errorf("reservation select balance failed: %w", err)
	}

	if currentBalance < model.Amount.Kopecks() {
		return models.ErrInsufficientFunds
	}

	hasConflict, checkErr := db.hasReservationConflict(ctx, tx, model)

	if checkErr != nil {
		return checkErr
	}

	if hasConflict {
		return models.ErrConflictReservation
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO reservation (user_id, service_id, order_id, amount, status)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (service_id, order_id) 
		DO UPDATE SET 
			amount = reservation.amount + EXCLUDED.amount,
			status = $5
	`, model.UserID, model.ServiceID, model.OrderID,
		model.Amount.Kopecks(), models.ReservationStatusReserved)

	if err != nil {
		return fmt.Errorf("reservation creation failed: %w", err)
	}

	_, err = tx.Exec(ctx, `
		UPDATE wallet 
		SET balance = balance - $1 
		WHERE user_id = $2
		`, model.Amount.Kopecks(), model.UserID)

	if err != nil {
		return fmt.Errorf("reservation balance update failed: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("couldn't commit reservation transaction: %w", err)
	}

	return nil
}

func (db *Repo) hasReservationConflict(ctx context.Context, tx pgx.Tx, reserve models.Reservation) (bool, error) {
	var existingUserID int
	err := tx.QueryRow(ctx, `
		SELECT user_id 
		FROM reservation 
		WHERE service_id = $1 AND order_id = $2
	`, reserve.ServiceID, reserve.OrderID).Scan(&existingUserID)

	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to check reservation conflict: %w", err)
	}

	return existingUserID != reserve.UserID, nil
}
