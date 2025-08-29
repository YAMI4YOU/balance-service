package models

import "math"

type Money struct {
	kopecks int64
}

func NewMoneyFromKopecks(kopecks int64) Money {
	return Money{kopecks: kopecks}
}

func NewMoneyFromRubles(rubles float64) Money {
	kopecks := int64(math.Round(rubles * 100))
	return Money{kopecks: kopecks}
}

func (m Money) Kopecks() int64 {
	return m.kopecks
}

func (m Money) Rubles() float64 {
	return float64(m.kopecks) / 100.0
}
