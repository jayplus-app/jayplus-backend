package app

import (
	"backend/contracts/db"
	"net/http"
)

type AppInterface interface {
	UICOnfig(w http.ResponseWriter, r *http.Request, db db.DBInterface)
	BookingConfig(w http.ResponseWriter, r *http.Request, db db.DBInterface)
}
