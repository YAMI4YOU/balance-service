package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Reserve struct {
	UserID    int   `json:"user_id"`
	ServiceID int   `json:"service_id"`
	OrderID   int   `json:"order_id"`
	Amount    int64 `json:"amount"`
}

func (db *DB) ReserveFunds(ctx context.Context, req Reserve) error {
	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("couldn't start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var currentBalance int64
	err = tx.QueryRow(ctx, "SELECT balance FROM wallet WHERE user_id = $1 FOR UPDATE",
		req.UserID).Scan(&currentBalance)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return fmt.Errorf("user not found")
	case err != nil:
		return fmt.Errorf("couldn't get current balance: %w", err)
	case currentBalance < req.Amount:
		return fmt.Errorf("insufficient balance")
	}

	_, err = tx.Exec(ctx, "UPDATE wallet SET balance = balance - $1 WHERE user_id = $2;",
		req.Amount, req.UserID)
	if err != nil {
		return fmt.Errorf("balance update failed: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("couldn't commit transaction: %w", err)
	}

	return nil
}
