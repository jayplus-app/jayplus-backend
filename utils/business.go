package utils

import (
	"backend/models"
	"time"
)

// IsWithinBusinessHours checks if a given time is within the specified business hours
func IsWithinBusinessHours(t time.Time, hours *models.BusinessHours) bool {
	var dayHours *models.DayHours

	switch t.Weekday() {
	case time.Monday:
		dayHours = &hours.Monday
	case time.Tuesday:
		dayHours = &hours.Tuesday
	case time.Wednesday:
		dayHours = &hours.Wednesday
	case time.Thursday:
		dayHours = &hours.Thursday
	case time.Friday:
		dayHours = &hours.Friday
	case time.Saturday:
		dayHours = &hours.Saturday
	case time.Sunday:
		dayHours = &hours.Sunday
	}

	if dayHours.Closed {
		return false
	}

	startTime, _ := time.Parse("15:04", dayHours.Start)
	endTime, _ := time.Parse("15:04", dayHours.End)

	// Convert the parsed times to the current day
	startTime = time.Date(t.Year(), t.Month(), t.Day(), startTime.Hour(), startTime.Minute(), 0, 0, t.Location())
	endTime = time.Date(t.Year(), t.Month(), t.Day(), endTime.Hour(), endTime.Minute(), 0, 0, t.Location())

	return t.After(startTime) && t.Before(endTime)
}
