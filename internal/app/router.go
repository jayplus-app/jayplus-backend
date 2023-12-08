package app

import (
	"backend/internal/auth"
	"backend/internal/booking"
	"backend/internal/payment"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func (app *App) SetupRouter() *mux.Router {
	r := mux.NewRouter()

	AppRoutes(r, app, app.DB)
	auth.AuthRoutes(r, app.Auth, app.DB)
	booking.BookingRoutes(r, app.Auth, app.DB)
	payment.PaymentRoutes(r, app.DB)

	corsOpts := handlers.AllowedOrigins([]string{"http://localhost:3000", "http://business1.localhost:3000", "https://jayplus.app", "https://business1.jayplus.app"})
    corsHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Business-Name"})
    corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	r.Use(handlers.CORS(corsOpts, corsHeaders, corsMethods))
	
	return r
}
