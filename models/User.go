package models

import "time"

type User struct {
	ID             int       `json:"id"`
	Email          string    `json:"email"`
	PhoneNumber    string    `json:"phoneNumber"`
	HashedPassword string    `json:"hashedPassword"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}
