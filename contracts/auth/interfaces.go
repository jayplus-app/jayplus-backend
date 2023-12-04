package auth

import (
	"backend/contracts/db"
	"net/http"

	"github.com/gorilla/mux"
)

type AuthInterface interface {
	// Handlers
	Login(w http.ResponseWriter, r *http.Request, db db.DBInterface)
	RefreshToken(w http.ResponseWriter, r *http.Request, db db.DBInterface)
	Logout(w http.ResponseWriter, r *http.Request)

	// Middleware
	AuthRequired(db db.DBInterface) mux.MiddlewareFunc
}
