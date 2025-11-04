package deposit

import (
	"context"

	"github.com/YAMI4YOU/balance-service/internal/models"
)

type repo interface {
	MakeDeposit(ctx context.Context, deposit models.Deposit) error
}
