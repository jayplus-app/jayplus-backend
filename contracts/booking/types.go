package booking

import "backend/models"

type Booking struct {
	Number string `json:"number"`
}

type VehicleTypes struct {
	Name  string                `json:"name"`
	Types []*models.VehicleType `json:"types"`
}

type ServiceTypes struct {
	Name  string                `json:"name"`
	Types []*models.ServiceType `json:"types"`
}

type TimeSlots struct {
	Date  string             `json:"date"`
	Slots []*models.TimeSlot `json:"slots"`
}

type Bookings struct {
	Date     string            `json:"date"`
	Bookings []*models.Booking `json:"bookings"`
}
