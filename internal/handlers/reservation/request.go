package reservation

import "errors"

type request struct {
	UserID    int     `json:"user_id"`
	ServiceID int     `json:"service_id"`
	OrderID   int     `json:"order_id"`
	Amount    float64 `json:"amount"`
}

func (req *request) Validate() error {
	if req.UserID <= 0 {
		return errors.New("invalid user_id")
	}
	if req.ServiceID <= 0 {
		return errors.New("invalid service_id")
	}
	if req.OrderID <= 0 {
		return errors.New("invalid order_id")
	}
	if req.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	return nil
}
