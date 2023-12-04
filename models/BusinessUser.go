package models

import "time"

type BusinessUser struct {
	BusinessID int       `json:"businessId"`
	UserID     int       `json:"userId"`
	RoleID     int       `json:"roleId"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	CreatedAt  time.Time `json:"createdAt"`
}
