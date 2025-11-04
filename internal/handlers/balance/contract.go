package balance

import (
	"context"

	"github.com/YAMI4YOU/balance-service/internal/models"
)

type repo interface {
	FetchBalance(ctx context.Context, userID int) (models.Money, error)
}
