package payment

import (
	"backend/contracts/db"
	"net/http"

	"github.com/gorilla/mux"
)

func PaymentRoutes(r *mux.Router, db db.DBInterface) {
	paymentRouter := r.PathPrefix("/payment").Subrouter()

	paymentRouter.HandleFunc("/create-payment-intent", func(w http.ResponseWriter, r *http.Request) {
		CreatePaymentIntent(w, r, db)
	}).Methods("POST", "OPTIONS")
	paymentRouter.HandleFunc("/booking-receipt/{bookingID}", func(w http.ResponseWriter, r *http.Request) {
		BookingReceipt(w, r, db)
	}).Methods("GET")
	paymentRouter.HandleFunc("/stripe-webhook", func(w http.ResponseWriter, r *http.Request) {
		StripeWebhook(w, r, db)
	}).Methods("POST")

	paymentRouter.HandleFunc("/pay-booking", func(w http.ResponseWriter, r *http.Request) {
		PayBooking(w, r, db)
	}).Methods("GET")
	paymentRouter.HandleFunc("/invoice/{booking-id}", func(w http.ResponseWriter, r *http.Request) {
		GetInvoice(w, r, db)
	}).Methods("GET")
}
