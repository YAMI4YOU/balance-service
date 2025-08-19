package db

import (
	"context"
	"fmt"
)

type DepositRequest struct {
	UserID  int   `json:"user_id"`
	Balance int64 `json:"balance"`
}

func (db *DB) Deposit(ctx context.Context, req DepositRequest) error {
	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("couldn't start a transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	/*var exists bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM wallet WHERE user_id = $1)", req.UserID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("couldn't check if user exists: %w", err)
	}
	fmt.Println("User exists:", exists)

	fmt.Println(req.UserID, req.Balance)
	if !exists {
		_, err := tx.Exec(ctx, "INSERT INTO wallet (user_id, balance) VALUES ($1, $2)", req.UserID, req.Balance)
		if err != nil {
			return fmt.Errorf("couldn't insert balance: %w", err)
		}
	} else {
		_, err := tx.Exec(ctx, "UPDATE wallet SET balance = balance + $1 WHERE user_id = $2", req.Balance, req.UserID)
		if err != nil {
			return fmt.Errorf("couldn't update balance: %w", err)
		}
	}*/

	_, err = tx.Exec(ctx, `
    INSERT INTO wallet (user_id, balance)
    VALUES ($1, $2)
    ON CONFLICT (user_id) 
    DO UPDATE SET balance = wallet.balance + EXCLUDED.balance
`, req.UserID, req.Balance)

	if err != nil {
		return fmt.Errorf("upsert wallet balance: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
