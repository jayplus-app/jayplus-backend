package models

type Transaction struct {
	ID        int    `json:"id"`
	BookingID int    `json:"bookingID"`
	Amount    int    `json:"amount"`
	Gateway   string `json:"gateway"`
	Status    string `json:"status"`
	Notes     string `json:"notes"`
}
