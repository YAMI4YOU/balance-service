package models

import "errors"

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrInsufficientFunds      = errors.New("insufficient funds")
	ErrReservationNotFound    = errors.New("reservation not found")
	ErrReservationNotReserved = errors.New("reservation is not in reserved status")
	ErrConflictReservation    = errors.New("order already reserved by different user")
	ErrUserMismatch           = errors.New("user mismatch")
	ErrAmountMismatch         = errors.New("amount mismatch")
)
