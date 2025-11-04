package reservation

import "github.com/YAMI4YOU/balance-service/internal/models"

type repo interface {
	ReserveFunds(ctx context.Context, req models.Reservation) error
}
