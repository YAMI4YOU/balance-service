package revenue

import "github.com/YAMI4YOU/balance-service/internal/models"

type repo interface {
	RecognizeRevenue(ctx context.Context, revenue models.Revenue) error
}
