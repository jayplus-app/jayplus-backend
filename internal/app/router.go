package app

import (
	"backend/internal/auth"
	"backend/internal/booking"
	"backend/internal/payment"

	"github.com/gorilla/mux"
)

func (app *App) SetupRouter() *mux.Router {
	r := mux.NewRouter()
	AppRoutes(r, app, app.DB)
	auth.AuthRoutes(r, app.Auth, app.DB)
	booking.BookingRoutes(r, app.Auth, app.DB)
	payment.PaymentRoutes(r, app.DB)
	return r
}
