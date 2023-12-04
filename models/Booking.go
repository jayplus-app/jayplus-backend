package models

import "time"

type Booking struct {
	ID               int       `json:"id"`
	BusinessID       int       `json:"businessId"`
	UserID           int       `json:"userId"`
	VehicleTypeID    int       `json:"vehicleTypeId"`
	ServiceTypeID    int       `json:"serviceTypeId"`
	Datetime         time.Time `json:"datetime"`
	EstimatedMinutes int       `json:"estimatedMinutes"`
	Cost             int       `json:"cost"`
	Discount         int       `json:"discount"`
	Deposit          int       `json:"deposit"`
	BillNumber       int       `json:"billNumber"`
	Status           string    `json:"status"`
}
