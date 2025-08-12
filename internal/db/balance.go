package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (db *DB) GetBalance(ctx context.Context, userID int) (int64, error) {
	const query = `SELECT balance FROM wallet WHERE user_id = $1`

	var balance int64
	err := db.Conn.QueryRow(ctx, query, userID).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("user not found")
		}
		return 0, fmt.Errorf("database query error: %w", err)
	}

	return balance, nil
}
