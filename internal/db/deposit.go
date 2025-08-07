package db

import (
	"context"
	"fmt"
)

type DepositRequest struct {
	UserID int   `json:"user_id"`
	Amount int64 `json:"amount"`
}

func (db *DB) Deposit(ctx context.Context, req DepositRequest) error {
	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("couldn't start a transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	var exists bool
	err = tx.QueryRow(ctx,
		"select exists(select 1 from balance where user_id = $1)", req.UserID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("couldn't check if user exists: %v", err)
	}

	if !exists {
		_, err := tx.Exec(ctx, "INSERT INTO balance (user_id, balance) VALUES ($1, $2)", req.UserID, req.Amount)
		if err != nil {
			return fmt.Errorf("couldn't insert balance: %v", err)
		}
	} else {
		_, err := tx.Exec(ctx, "UPDATE balance SET balance = balance + $1 WHERE user_id = $2", req.Amount, req.UserID)
		if err != nil {
			return fmt.Errorf("couldn't update balance: %v", err)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("couldn't commit transaction: %v", err)
		}
	}

	return nil
}
