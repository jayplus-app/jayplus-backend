package payment

import (
	"backend/config"
	"backend/contracts/db"
	"backend/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/webhook"
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
	stripe.Key = config.StripeSecretKey

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

	businessName := r.Header.Get("Business-Name")

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
		Metadata: map[string]string{
			"bookingID":    strconv.Itoa(booking.ID),
			"businessName": businessName,
		},
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

	bookingID, err := strconv.Atoi(vars["bookingID"])
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

func StripeWebhook(w http.ResponseWriter, r *http.Request, db db.DBInterface) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	event := stripe.Event{}

	if err := json.Unmarshal(payload, &event); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Webhook error while parsing basic request. %v\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	endpointSecret := config.StripeWebhookSecret
	signatureHeader := r.Header.Get("Stripe-Signature")

	event, err = webhook.ConstructEvent(payload, signatureHeader, endpointSecret)
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Webhook signature verification failed. %v\n", err)
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}

	// Unmarshal the event data into an appropriate struct depending on its Type
	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = PaymentIntentSucceeded(db, paymentIntent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error handling webhook: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	}

	w.WriteHeader(http.StatusOK)
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
