package models

type BookingLog struct {
	ID        int    `json:"id"`
	BookingID int    `json:"bookingId"`
	UserID    int    `json:"userId"`
	State     string `json:"state"`
	Details   string `json:"details"`
}
