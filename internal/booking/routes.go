package booking

import (
	"backend/contracts/auth"
	"backend/contracts/db"
	"net/http"

	"github.com/gorilla/mux"
)

func BookingRoutes(r *mux.Router, auth auth.AuthInterface, db db.DBInterface) {
	bookingRouter := r.PathPrefix("/booking").Subrouter()

	bookingRouter.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		Test(w, r, db)
	}).Methods("GET")

	bookingRouter.HandleFunc("/vehicle-types", func(w http.ResponseWriter, r *http.Request) {
		VehicleTypes(w, r, db)
	}).Methods("GET")
	bookingRouter.HandleFunc("/service-types", func(w http.ResponseWriter, r *http.Request) {
		ServiceTypes(w, r, db)
	}).Methods("GET")
	bookingRouter.HandleFunc("/timeslots", func(w http.ResponseWriter, r *http.Request) {
		TimeSlots(w, r, db)
	}).Methods("POST")
	bookingRouter.HandleFunc("/service-cost", func(w http.ResponseWriter, r *http.Request) {
		ServiceCost(w, r, db)
	}).Methods("POST")
	bookingRouter.HandleFunc("/create-booking", func(w http.ResponseWriter, r *http.Request) {
		CreateBooking(w, r, db)
	}).Methods("POST")

	adminOnlyRouter := bookingRouter.PathPrefix("/").Subrouter()
	adminOnlyRouter.Use(auth.AuthRequired(db))

	adminOnlyRouter.HandleFunc("/bookings", func(w http.ResponseWriter, r *http.Request) {
		Bookings(w, r, db)
	}).Methods("GET")
	adminOnlyRouter.HandleFunc("/booking/{id}", func(w http.ResponseWriter, r *http.Request) {
		Booking(w, r, db)
	}).Methods("GET")
	adminOnlyRouter.HandleFunc("/create-booking-admin", func(w http.ResponseWriter, r *http.Request) {
		CreateBookingAdmin(w, r, db)
	}).Methods("POST")
	adminOnlyRouter.HandleFunc("/cancel-booking/{id}", func(w http.ResponseWriter, r *http.Request) {
		CancelBooking(w, r, db)
	}).Methods("GET")
}
