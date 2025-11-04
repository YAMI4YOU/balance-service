package models

type ReservationStatus string

const (
	ReservationStatusReserved  ReservationStatus = "reserved"
	ReservationStatusConfirmed ReservationStatus = "confirmed"
	ReservationStatusCancelled ReservationStatus = "cancelled"
)

type Reservation struct {
	UserID    int
	ServiceID int
	OrderID   int
	Amount    Money
}
