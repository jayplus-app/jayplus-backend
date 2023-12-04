package db

import (
	"backend/models"
	"time"
)

type DBInterface interface {
	SetupDB() error

	// App
	GetBusinessUIConfigByID(business_id int) (*models.UIConfig, error)
	GetBusinessBookingConfigByID(business_id int) (*models.BookingConfig, error)

	// Common
	GetBusinessByBusinessName(businessName string) (*models.Business, error)

	// Auth
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	IsUserInBusiness(userID, businessID int) (bool, error)
	GetRoleByBusinessIDAndUserID(businessID, userID int) (*models.Role, error)
	GetBusinessHoursByBusinessID(businessID int) (*models.BusinessHours, error)

	// Booking
	GetVehicleTypes(businessID int) ([]*models.VehicleType, error)
	GetVehicleTypeByID(vehicleTypeID int) (*models.VehicleType, error)
	GetServiceTypes(businessID int) ([]*models.ServiceType, error)
	GetServiceTypeByID(serviceTypeID int) (*models.ServiceType, error)
	GetServiceDetail(businessID, vehicleTypeID, serviceTypeID int) (*models.ServiceDetail, error)
	GetBookingsByDate(businessID int, date time.Time) ([]*models.Booking, error)
	GetBookingByID(bookingID int) (*models.Booking, error)
	CreateBooking(booking *models.Booking) (*models.Booking, error)
	CreateBookingLog(bookingLog *models.BookingLog) error
	UpdateBookingStatus(bookingID int, status string) error

	// Payment
	RecordPayment(payment *models.Payment) error
	GetTransactionByBookingID(transactionID int) (*models.Transaction, error)

	// SMS
	RecordSMS(sms *models.SMS) error
}
