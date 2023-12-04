package auth

import (
	"net/http"

	"backend/contracts/auth"
	"backend/contracts/db"

	"github.com/gorilla/mux"
)

func AuthRoutes(r *mux.Router, auth auth.AuthInterface, db db.DBInterface) {
	authRouter := r.PathPrefix("/auth").Subrouter()

	authRouter.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		auth.Login(w, r, db)
	}).Methods("POST")
	authRouter.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		auth.RefreshToken(w, r, db)
	}).Methods("GET")
	authRouter.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		auth.Logout(w, r)
	}).Methods("GET")
}
