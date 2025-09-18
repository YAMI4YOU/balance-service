package models

type Report struct {
	Year  int
	Month int
}

type ReportSummary struct {
	ServiceID int64
	Amount    Money
	CreatedAt string
}
