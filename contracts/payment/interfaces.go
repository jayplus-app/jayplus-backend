package payment

import (
	"backend/contracts/db"
	"net/http"
)

type PaymentInterface interface {
	PayBooking(w http.ResponseWriter, r *http.Request, db db.DBInterface)
	GetInvoice(w http.ResponseWriter, r *http.Request, db db.DBInterface)
}
