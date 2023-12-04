package models

import "time"

type TimeSlot struct {
	ID          string    `json:"id"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	FreeMinutes int       `json:"freeMinutes"`
	Available   bool      `json:"available"`
	IsPast      bool      `json:"isPast"`
}
