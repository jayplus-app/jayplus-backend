package booking

import (
	"backend/contracts/db"
	"backend/models"
	"time"
)

// GetBookingTimeslots retrieves all booking timeslots.
func GetBookingTimeslots(db db.DBInterface, business *models.Business, serviceTypeID int, vehicleTypeID int, date time.Time) ([]*models.TimeSlot, error) {
	businessTZ, err := time.LoadLocation(business.Timezone)
	if err != nil {
		return nil, err
	}

	// TODO: Rename GetServiceDetail to GetServiceDetails
	serviceDetails, err := db.GetServiceDetail(business.ID, vehicleTypeID, serviceTypeID)
	if err != nil {
		return nil, err
	}

	// Initialize primitives
	capacity := 3         // TODO: [THREAD:3] Read from config
	timeslotMinutes := 60 // TODO: [THREAD:3] Read from config
	timeslotManminutes := timeslotMinutes * capacity
	threshold := float32(50) / 100 // TODO: [THREAD:3] Read from config
	maxOverflow := int(float32(timeslotManminutes) * threshold)

	// Set up the reception start and end times
	startHour := 9   // TODO: [THREAD:3] Read from config
	startMinute := 0 // TODO: [THREAD:3] Read from config
	endHour := 17    // TODO: [THREAD:3] Read from config
	endMinute := 30  // TODO: [THREAD:3] Read from config
	receptionStart := time.Date(date.Year(), date.Month(), date.Day(), startHour, startMinute, 0, 0, businessTZ)
	receptionEnd := time.Date(date.Year(), date.Month(), date.Day(), endHour, endMinute, 0, 0, businessTZ)
	start := receptionStart
	end := start.Add(time.Duration(timeslotMinutes) * time.Minute)

	// Retrieve booked timeslots within the day
	bookings, err := db.GetBookingsByDate(business.ID, date)
	if err != nil {
		return nil, err
	}

	// Generate timeslots
	timeslotCount := int(receptionEnd.Sub(receptionStart).Minutes()) / timeslotMinutes
	timeslots := make([]*models.TimeSlot, 0, timeslotCount)

	lastOverflow := 0

	for i := 0; i < timeslotCount; i++ {
		isPast := time.Now().After(start)
		isLastTimeslot := receptionEnd.Before(end)

		sum := 0
		sumNext := 0
		for _, b := range bookings {
			if b.Datetime.Equal(start) {
				sum += b.EstimatedMinutes
			} else if b.Datetime.Equal(end) {
				sumNext += b.EstimatedMinutes
			}
		}

		nextOverflow := 0
		if sumNext > timeslotManminutes {
			nextOverflow = sumNext - timeslotManminutes
		}

		allowedOverflow := 0
		if !isLastTimeslot && maxOverflow > nextOverflow {
			allowedOverflow = maxOverflow - nextOverflow
		}

		maxCalc := timeslotManminutes + allowedOverflow - lastOverflow
		remained := maxCalc - sum
		available := !isPast && (remained+allowedOverflow-serviceDetails.DurationMinutes >= 0)

		timeslot := models.TimeSlot{
			StartTime:   start,
			EndTime:     end,
			FreeMinutes: remained,
			Available:   available,
			IsPast:      isPast,
		}

		timeslots = append(timeslots, &timeslot)

		start = end
		end = start.Add(time.Duration(timeslotMinutes) * time.Minute)

		// Update last overflow
		lastOverflow = sum - timeslotManminutes - lastOverflow
		if lastOverflow < 0 {
			lastOverflow = 0
		}
	}

	return timeslots, nil
}
