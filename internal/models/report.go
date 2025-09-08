package models

type MonthlyReportRequest struct {
	Year  int
	Month int
}

type MonthlyReportRecord struct {
	UserID    int64
	ServiceID int64
	OrderID   int64
	Amount    Money
	Operation string
	CreatedAt string
}
