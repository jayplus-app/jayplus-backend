package booking

import (
	"backend/contracts/auth"
	"backend/contracts/booking"
	"backend/contracts/db"
	"backend/models"
	"backend/utils"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// VehicleTypes handler returns a list of vehicle types.
func Test(w http.ResponseWriter, r *http.Request, db db.DBInterface) {
	// test json response
	type ResponseBody struct {
		Message string `json:"message"`
	}

	responseBody := ResponseBody{
		Message: "Hello World",
	}

	utils.WriteJSON(w, http.StatusOK, responseBody)
}

// VehicleTypes handler returns a list of vehicle types.
func VehicleTypes(w http.ResponseWriter, r *http.Request, db db.DBInterface) {
	businessName := r.Header.Get("Business-Name")

	log.Println(businessName)

	business, err := db.GetBusinessByBusinessName(businessName)
	if err != nil {
		utils.ErrorJSON(w, errors.New("invalid business"), http.StatusBadRequest)
		return
	}

	vehicleTypes, err := db.GetVehicleTypes(business.ID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	utils.WriteJSON(w, http.StatusOK, vehicleTypes)
}

// ServiceTypes handler returns a list of service types.
func ServiceTypes(w http.ResponseWriter, r *http.Request, db db.DBInterface) {
	businessName := r.Header.Get("Business-Name")

	business, err := db.GetBusinessByBusinessName(businessName)
	if err != nil {
		utils.ErrorJSON(w, errors.New("invalid business"), http.StatusBadRequest)
		return
	}

	serviceTypes, err := db.GetServiceTypes(business.ID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	utils.WriteJSON(w, http.StatusOK, serviceTypes)
}

// Timeslots handler returns a list of time slots.
func TimeSlots(w http.ResponseWriter, r *http.Request, db db.DBInterface) {
	type RequestBody struct {
		StartDate     string `json:"dateTime"`
		ServiceTypeID string `json:"serviceTypeID"`
		VehicleTypeID string `json:"vehicleTypeID"`
	}

	var requestBody RequestBody
	if err := utils.ReadJSON(w, r, &requestBody); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	businessName := r.Header.Get("Business-Name")
	business, err := db.GetBusinessByBusinessName(businessName)
	if err != nil {
		// TODO: log error
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	vehicleTypeID, err := strconv.Atoi(requestBody.VehicleTypeID)
	if err != nil {
		// TODO: log error
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	serviceTypeID, err := strconv.Atoi(requestBody.ServiceTypeID)
	if err != nil {
		// TODO: log error
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	d, err := time.Parse(time.DateOnly, requestBody.StartDate)
	if err != nil {
		// TODO: log error
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	timeSlots, err := GetBookingTimeslots(db, business, serviceTypeID, vehicleTypeID, d)
	if err != nil {
		// TODO: log error
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	timeSlotsResponse := booking.TimeSlots{
		Date:  d.Format(time.DateOnly),
		Slots: timeSlots,
	}

	utils.WriteJSON(w, http.StatusOK, timeSlotsResponse)
}

// ServiceCost handler returns the cost of a service.
func ServiceCost(w http.ResponseWriter, r *http.Request, db db.DBInterface) {
	// post request body
	type RequestBody struct {
		ServiceTypeID string `json:"serviceTypeID"`
		VehicleTypeID string `json:"vehicleTypeID"`
	}

	var requestBody RequestBody
	err := utils.ReadJSON(w, r, &requestBody)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	businessName := r.Header.Get("Business-Name")
	business, err := db.GetBusinessByBusinessName(businessName)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	vehicleTypeID, err := strconv.Atoi(requestBody.VehicleTypeID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	serviceTypeID, err := strconv.Atoi(requestBody.ServiceTypeID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	serviceDetail, err := db.GetServiceDetail(business.ID, vehicleTypeID, serviceTypeID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	utils.WriteJSON(w, http.StatusOK, serviceDetail)
}

type BookingsResponse struct {
	ID          int       `json:"id"`
	VehicleType string    `json:"vehicleType"`
	ServiceType string    `json:"serviceType"`
	Datetime    time.Time `json:"datetime"`
	Status      string    `json:"status"`
}

// Bookings handler returns a list of bookings.
func Bookings(w http.ResponseWriter, r *http.Request, db db.DBInterface) {
	dateInput := r.URL.Query().Get("date")
	if dateInput == "" {
		utils.ErrorJSON(w, errors.New("invalid date"), http.StatusBadRequest)
		return
	}

	businessName := r.Header.Get("Business-Name")
	business, err := db.GetBusinessByBusinessName(businessName)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	businessTZ, err := time.LoadLocation(business.Timezone)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", dateInput)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	var bookingsResponse []BookingsResponse

	bookings, err := db.GetBookingsByDate(business.ID, date)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	for _, booking := range bookings {
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

		bookingResponse := BookingsResponse{
			ID:          booking.ID,
			VehicleType: vehicleType.Name,
			ServiceType: serviceType.Name,
			Datetime:    booking.Datetime.In(businessTZ),
			Status:      booking.Status,
		}

		bookingsResponse = append(bookingsResponse, bookingResponse)
	}

	utils.WriteJSON(w, http.StatusOK, bookingsResponse)
}

type BookingResponse struct {
	ID          int       `json:"id"`
	UserID      int       `json:"userID"`
	VehicleType string    `json:"vehicleType"`
	ServiceType string    `json:"serviceType"`
	Datetime    time.Time `json:"datetime"`
	Cost        int       `json:"cost"`
	Discount    int       `json:"discount"`
	Deposit     int       `json:"deposit"`
	BillNumber  int       `json:"billNumber"`
	Status      string    `json:"status"`
}

// Booking handler returns a booking.
func Booking(w http.ResponseWriter, r *http.Request, db db.DBInterface) {
	vars := mux.Vars(r)

	bookingID, err := strconv.Atoi(vars["id"])
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

	bookingResponse := BookingResponse{
		ID:          booking.ID,
		UserID:      booking.UserID,
		VehicleType: vehicleType.Name,
		ServiceType: serviceType.Name,
		Datetime:    booking.Datetime,
		Cost:        booking.Cost,
		Discount:    booking.Discount,
		Deposit:     booking.Deposit,
		BillNumber:  booking.BillNumber,
		Status:      booking.Status,
	}

	utils.WriteJSON(w, http.StatusOK, bookingResponse)
}

type CancelBookingResponse struct {
	Message string `json:"message"`
}

// CancelBooking handler cancels a booking.
func CancelBooking(w http.ResponseWriter, r *http.Request, db db.DBInterface) {
	bookingID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	booking, err := db.GetBookingByID(bookingID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	if booking.Status == "cancelled" {
		utils.ErrorJSON(w, errors.New("booking already cancelled"), http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value("claims").((*auth.JWTClaims))
	if !ok {
		utils.ErrorJSON(w, errors.New("failed to retrieve claims"), http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	businessID := claims.BusinessID

	role, err := db.GetRoleByBusinessIDAndUserID(businessID, userID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	err = db.UpdateBookingStatus(bookingID, "cancelled")
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	bookingLog := models.BookingLog{
		BookingID: bookingID,
		UserID:    userID,
		State:     "cancelled",
		Details:   "Booking cancelled by " + role.Name + " with UserID " + strconv.Itoa(userID),
	}

	err = db.CreateBookingLog(&bookingLog)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	cancelBookingResponse := CancelBookingResponse{
		Message: "Booking cancelled successfully",
	}

	utils.WriteJSON(w, http.StatusOK, cancelBookingResponse)
}

type CreateBookingResponse struct {
	ID               int       `json:"id"`
	UserID           int       `json:"userID"`
	VehicleType      string    `json:"vehicleType"`
	ServiceType      string    `json:"serviceType"`
	Datetime         time.Time `json:"datetime"`
	EstimatedMinutes int       `json:"estimatedMinutes"`
	Cost             int       `json:"cost"`
	Discount         int       `json:"discount"`
	Deposit          int       `json:"deposit"`
	BillNumber       int       `json:"billNumber"`
	Status           string    `json:"status"`
}

// CreateBooking handler creates bookings.
func CreateBooking(w http.ResponseWriter, r *http.Request, db db.DBInterface) {
	type RequestBody struct {
		VehicleTypeID string `json:"vehicleTypeID"`
		ServiceTypeID string `json:"serviceTypeID"`
		Datetime      string `json:"datetime"`
	}

	var requestBody RequestBody
	err := utils.ReadJSON(w, r, &requestBody)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	businessName := r.Header.Get("Business-Name")
	business, err := db.GetBusinessByBusinessName(businessName)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	userID := 2

	vehicleTypeID, err := strconv.Atoi(requestBody.VehicleTypeID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	serviceTypeID, err := strconv.Atoi(requestBody.ServiceTypeID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	datetime, err := time.Parse(time.RFC3339, requestBody.Datetime)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Check if booking is in the past
	if datetime.Before(time.Now()) {
		utils.ErrorJSON(w, errors.New("booking is in the past"), http.StatusBadRequest)
		return
	}

	// Check if booking is within business hours
	businessHours, err := db.GetBusinessHoursByBusinessID(business.ID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Check if booking is within business hours
	if !utils.IsWithinBusinessHours(datetime, businessHours) {
		utils.ErrorJSON(w, errors.New("booking is not within business hours"), http.StatusBadRequest)
		return
	}

	serviceDetail, err := db.GetServiceDetail(business.ID, vehicleTypeID, serviceTypeID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	discount := 0

	deposit := serviceDetail.Price * 20 / 100

	booking := models.Booking{
		BusinessID:       business.ID,
		UserID:           userID,
		VehicleTypeID:    vehicleTypeID,
		ServiceTypeID:    serviceTypeID,
		Datetime:         datetime,
		EstimatedMinutes: serviceDetail.DurationMinutes,
		Cost:             serviceDetail.Price,
		Discount:         discount,
		Deposit:          deposit,
		BillNumber:       0,
		Status:           "pending",
	}

	returnedBooking, err := db.CreateBooking(&booking)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	bookingLog := models.BookingLog{
		BookingID: returnedBooking.ID,
		UserID:    userID,
		State:     "created",
		Details:   "Booking created by " + "unknown" + " with UserID " + strconv.Itoa(userID),
	}

	err = db.CreateBookingLog(&bookingLog)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	vehicleType, err := db.GetVehicleTypeByID(returnedBooking.VehicleTypeID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	serviceType, err := db.GetServiceTypeByID(returnedBooking.ServiceTypeID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	bookingResponse := CreateBookingResponse{
		ID:               returnedBooking.ID,
		UserID:           returnedBooking.UserID,
		VehicleType:      vehicleType.Name,
		ServiceType:      serviceType.Name,
		Datetime:         returnedBooking.Datetime,
		EstimatedMinutes: returnedBooking.EstimatedMinutes,
		Cost:             returnedBooking.Cost,
		Discount:         returnedBooking.Discount,
		Deposit:          returnedBooking.Deposit,
		BillNumber:       returnedBooking.BillNumber,
		Status:           returnedBooking.Status,
	}

	utils.WriteJSON(w, http.StatusOK, bookingResponse)
}

// CreateBookingAdmin handler creates bookings.
func CreateBookingAdmin(w http.ResponseWriter, r *http.Request, db db.DBInterface) {
	type RequestBody struct {
		VehicleTypeID string `json:"vehicleTypeID"`
		ServiceTypeID string `json:"serviceTypeID"`
		Datetime      string `json:"datetime"`
	}

	var requestBody RequestBody
	err := utils.ReadJSON(w, r, &requestBody)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	businessName := r.Header.Get("Business-Name")
	business, err := db.GetBusinessByBusinessName(businessName)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	var userID int

	claims, ok := r.Context().Value("claims").((*auth.JWTClaims))
	if !ok {
		utils.ErrorJSON(w, errors.New("failed to retrieve claims"), http.StatusBadRequest)
		return
	}

	userID, err = strconv.Atoi(claims.Subject)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	vehicleTypeID, err := strconv.Atoi(requestBody.VehicleTypeID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	serviceTypeID, err := strconv.Atoi(requestBody.ServiceTypeID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	datetime, err := time.Parse(time.RFC3339, requestBody.Datetime)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Check if booking is in the past
	if datetime.Before(time.Now()) {
		utils.ErrorJSON(w, errors.New("booking is in the past"), http.StatusBadRequest)
		return
	}

	// Check if booking is within business hours
	businessHours, err := db.GetBusinessHoursByBusinessID(business.ID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Check if booking is within business hours
	if !utils.IsWithinBusinessHours(datetime, businessHours) {
		utils.ErrorJSON(w, errors.New("booking is not within business hours"), http.StatusBadRequest)
		return
	}

	serviceDetail, err := db.GetServiceDetail(business.ID, vehicleTypeID, serviceTypeID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	discount := 0

	deposit := 0

	booking := models.Booking{
		BusinessID:       business.ID,
		UserID:           userID,
		VehicleTypeID:    vehicleTypeID,
		ServiceTypeID:    serviceTypeID,
		Datetime:         datetime,
		EstimatedMinutes: serviceDetail.DurationMinutes,
		Cost:             serviceDetail.Price,
		Discount:         discount,
		Deposit:          deposit,
		BillNumber:       0,
		Status:           "active",
	}

	returnedBooking, err := db.CreateBooking(&booking)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	role, err := db.GetRoleByBusinessIDAndUserID(returnedBooking.BusinessID, userID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	bookingLog := models.BookingLog{
		BookingID: returnedBooking.ID,
		UserID:    userID,
		State:     "created",
		Details:   "Booking created by " + role.Name + " with UserID " + strconv.Itoa(userID),
	}

	err = db.CreateBookingLog(&bookingLog)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	vehicleType, err := db.GetVehicleTypeByID(returnedBooking.VehicleTypeID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	serviceType, err := db.GetServiceTypeByID(returnedBooking.ServiceTypeID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	bookingResponse := CreateBookingResponse{
		ID:               returnedBooking.ID,
		UserID:           returnedBooking.UserID,
		VehicleType:      vehicleType.Name,
		ServiceType:      serviceType.Name,
		Datetime:         returnedBooking.Datetime,
		EstimatedMinutes: returnedBooking.EstimatedMinutes,
		Cost:             returnedBooking.Cost,
		Discount:         returnedBooking.Discount,
		Deposit:          returnedBooking.Deposit,
		BillNumber:       returnedBooking.BillNumber,
		Status:           returnedBooking.Status,
	}

	utils.WriteJSON(w, http.StatusOK, bookingResponse)
}
