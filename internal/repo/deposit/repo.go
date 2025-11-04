package deposit

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/YAMI4YOU/balance-service/internal/models"
)

type Repo struct {
	conn *pgx.Conn
}

func New(conn *pgx.Conn) *Repo {
	return &Repo{conn: conn}
}

func (db *Repo) MakeDeposit(ctx context.Context, model models.Deposit) error {
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("couldn't start a transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
    INSERT INTO wallet (user_id, balance)
    VALUES ($1, $2)
    ON CONFLICT (user_id) 
    DO UPDATE SET balance = wallet.balance + EXCLUDED.balance
`, model.UserID, model.Balance.Kopecks())

	if err != nil {
		return fmt.Errorf("upsert wallet balance: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
