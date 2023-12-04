package app

import (
	"backend/contracts/app"
	"backend/contracts/db"
	"net/http"

	"github.com/gorilla/mux"
)

func AppRoutes(r *mux.Router, app app.AppInterface, db db.DBInterface) {
	appRouter := r.PathPrefix("/app").Subrouter()

	appRouter.HandleFunc("/ui-config", func(w http.ResponseWriter, r *http.Request) {
		app.UICOnfig(w, r, db)
	}).Methods("GET")
	appRouter.HandleFunc("/booking-config", func(w http.ResponseWriter, r *http.Request) {
		app.BookingConfig(w, r, db)
	}).Methods("GET")
}
