package booking

import (
	"backend/contracts/db"
	"net/http"
)

type BookingInterface interface {
	// Handlers
	Test(w http.ResponseWriter, r *http.Request, db db.DBInterface)
	VehicleTypes(w http.ResponseWriter, r *http.Request, db db.DBInterface)
	ServiceTypes(w http.ResponseWriter, r *http.Request, db db.DBInterface)
	TimeSlots(w http.ResponseWriter, r *http.Request, db db.DBInterface)
	ServiceCost(w http.ResponseWriter, r *http.Request, db db.DBInterface)
	CreateBooking(w http.ResponseWriter, r *http.Request, db db.DBInterface)
	CreateBookingAdmin(w http.ResponseWriter, r *http.Request, db db.DBInterface)
	Bookings(w http.ResponseWriter, r *http.Request, db db.DBInterface)
	Booking(w http.ResponseWriter, r *http.Request, db db.DBInterface)
	CancelBooking(w http.ResponseWriter, r *http.Request, db db.DBInterface)
}
