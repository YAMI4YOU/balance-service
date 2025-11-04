package models

type Revenue struct {
	UserID    int
	ServiceID int
	OrderID   int
	Amount    Money
}
