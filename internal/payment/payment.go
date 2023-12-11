package payment

import (
	"backend/contracts/db"
	"backend/models"
	"errors"
	"strconv"

	"github.com/stripe/stripe-go/v76"
)

func PaymentIntentSucceeded(db db.DBInterface, pi stripe.PaymentIntent) error {
	bookingID, err := strconv.Atoi(pi.Metadata["bookingID"])
	if err != nil {
		return errors.New("error converting booking id to int")
	}

	businessName := pi.Metadata["businessName"]

	business, err := db.GetBusinessByBusinessName(businessName)

	// TODO: get role from token. Number 2 is hardcoded for now since it represents the unknown user
	userID := 2

	err = db.UpdateBookingStatus(bookingID, "active")
	if err != nil {
		return errors.New("error updating booking status")
	}

	bookingLog := models.BookingLog{
		BookingID: bookingID,
		UserID:    userID,
		State:     "active",
		Details:   "Booking deposit paid by unknown user with UserID " + strconv.Itoa(userID) + " for business " + businessName + " with ID " + strconv.Itoa(business.ID) + ".",
	}

	err = db.CreateBookingLog(&bookingLog)
	if err != nil {
		return errors.New("error creating booking log")
	}

	return nil
}
