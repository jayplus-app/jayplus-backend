package payment

import (
	"backend/contracts/db"
	"backend/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func PayBooking(w http.ResponseWriter, r *http.Request, db db.DBInterface) {
	payment := map[string]string{
		"booking_id": "1",
		"amount":     "100",
	}

	utils.WriteJSON(w, http.StatusOK, payment)
}

type GetInvoiceResponse struct {
	BookingID     int       `json:"bookingID"`
	TransactionID int       `json:"transactionID"`
	BillNumber    int       `json:"billNumber"`
	Status        string    `json:"status"`
	VehicleType   string    `json:"vehicleType"`
	ServiceType   string    `json:"serviceType"`
	Datetime      time.Time `json:"datetime"`
	ServiceCost   int       `json:"serviceCost"`
	Discount      int       `json:"discount"`
	Total         int       `json:"total"`
	Deposit       int       `json:"deposit"`
	Remaining     int       `json:"remaining"`
}

// GetInvoice handler returns an invoice.
func GetInvoice(w http.ResponseWriter, r *http.Request, db db.DBInterface) {
	bookingID, err := strconv.Atoi(mux.Vars(r)["booking-id"])
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	booking, err := db.GetBookingByID(bookingID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	transaction, err := db.GetTransactionByBookingID(bookingID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	vehicleType, err := db.GetVehicleTypeByID(booking.VehicleTypeID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	serviceType, err := db.GetServiceTypeByID(booking.ServiceTypeID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	invoiceResponse := GetInvoiceResponse{
		BookingID:     booking.ID,
		TransactionID: transaction.ID,
		BillNumber:    booking.BillNumber,
		Status:        transaction.Status,
		VehicleType:   vehicleType.Name,
		ServiceType:   serviceType.Name,
		Datetime:      booking.Datetime,
		ServiceCost:   booking.Cost,
		Discount:      booking.Discount,
		Total:         booking.Cost - booking.Discount,
		Deposit:       booking.Deposit,
		Remaining:     booking.Cost - booking.Discount - booking.Deposit,
	}

	utils.WriteJSON(w, http.StatusOK, invoiceResponse)
}
