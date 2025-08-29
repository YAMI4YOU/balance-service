package balance

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/jackc/pgx/v5"

	"github.com/YAMI4YOU/balance-service/internal/models"
)

type Repo struct {
	conn *pgx.Conn
}

func New(conn *pgx.Conn) *Repo {
	return &Repo{conn: conn}
}

func (db *Repo) FetchBalance(ctx context.Context, userID int) (models.Money, error) {
	const query = `SELECT balance FROM wallet WHERE user_id = $1`

	var balance float64
	err := db.conn.QueryRow(ctx, query, userID).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.NewMoneyFromKopecks(0), models.ErrUserNotFound
		}
		return models.NewMoneyFromKopecks(0), fmt.Errorf("fetch balance error: %w", err)
	}
	kopecks := int64(math.Round(balance * 100))
	return models.NewMoneyFromKopecks(kopecks), nil
}
