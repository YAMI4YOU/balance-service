package report

import "errors"

const (
	startYear = 2000
	endYear   = 2100

	firstMonth = 1
	lastMonth  = 12
)

type request struct {
	Year  int `json:"year"`
	Month int `json:"month"`
}

func (req *request) Validate() error {
	if req.Year < startYear || req.Year > endYear {
		return errors.New("year out of range")
	}
	if req.Month < firstMonth || req.Month > lastMonth {
		return errors.New("month out of range")
	}

	return nil
}
