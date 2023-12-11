package payment

import (
	"backend/contracts/db"
	"backend/utils"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
)

type BookingInvoice struct {
	BookingID   int       `json:"bookingID"`
	Status      string    `json:"status"`
	VehicleType string    `json:"vehicleType"`
	ServiceType string    `json:"serviceType"`
	Datetime    time.Time `json:"datetime"`
	ServiceCost int       `json:"serviceCost"`
	Discount    int       `json:"discount"`
	Total       int       `json:"total"`
	Deposit     int       `json:"deposit"`
	Remaining   int       `json:"remaining"`
}

type CreatePaymentIntentResponse struct {
	BookingInvoice BookingInvoice `json:"bookingInvoice"`
	ClientSecret   string         `json:"clientSecret"`
}

func CreatePaymentIntent(w http.ResponseWriter, r *http.Request, db db.DBInterface) {
	stripe.Key = "sk_test_51MKWmeLKOrXppaRnyQt8qStlXASK4M5V7JOBQ2y3uTrJQWWWumXn5XEA0NnPTNFEEygMlx7ShKU2qqmeVg0lR1Rb00Md93C9ad"

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var reqPayload struct {
		BookingNumber int `json:"bookingNumber"`
	}

	err := utils.ReadJSON(w, r, &reqPayload)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	booking, err := db.GetBookingByID(reqPayload.BookingNumber)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Check if the booking is already paid.
	if booking.Status != "pending" {
		err := errors.New("booking is already paid")
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	serviceDetail, err := db.GetServiceDetail(booking.BusinessID, booking.VehicleTypeID, booking.ServiceTypeID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Create a PaymentIntent with the order amount and currency.
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(serviceDetail.Price)),
		Currency: stripe.String(string(stripe.CurrencyCAD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := paymentintent.New(params)
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

	bookingInvoice := BookingInvoice{
		BookingID:   booking.ID,
		Status:      booking.Status,
		VehicleType: vehicleType.Name,
		ServiceType: serviceType.Name,
		Datetime:    booking.Datetime,
		ServiceCost: booking.Cost,
		Discount:    booking.Discount,
		Total:       booking.Cost - booking.Discount,
		Deposit:     booking.Deposit,
		Remaining:   booking.Cost - booking.Discount - booking.Deposit,
	}

	response := CreatePaymentIntentResponse{
		BookingInvoice: bookingInvoice,
		ClientSecret:   pi.ClientSecret,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

type BookingReceiptResponse struct {
	BookingID   int       `json:"bookingID"`
	BillNumber  int       `json:"billNumber"`
	Status      string    `json:"status"`
	VehicleType string    `json:"vehicleType"`
	ServiceType string    `json:"serviceType"`
	Datetime    time.Time `json:"datetime"`
	ServiceCost int       `json:"serviceCost"`
	Discount    int       `json:"discount"`
	Total       int       `json:"total"`
	Deposit     int       `json:"deposit"`
	Remaining   int       `json:"remaining"`
}

func BookingReceipt(w http.ResponseWriter, r *http.Request, db db.DBInterface) {
	vars := mux.Vars(r)

	bookingID, err := strconv.Atoi(vars["booking-id"])
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	booking, err := db.GetBookingByID(bookingID)
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

	response := BookingReceiptResponse{
		BookingID:   booking.ID,
		BillNumber:  booking.BillNumber,
		Status:      booking.Status,
		VehicleType: vehicleType.Name,
		ServiceType: serviceType.Name,
		Datetime:    booking.Datetime,
		ServiceCost: booking.Cost,
		Discount:    booking.Discount,
		Total:       booking.Cost - booking.Discount,
		Deposit:     booking.Deposit,
		Remaining:   booking.Cost - booking.Discount - booking.Deposit,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

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
