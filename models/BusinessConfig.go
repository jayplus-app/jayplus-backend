package models

import "time"

type BusinessConfig struct {
	ID         int       `json:"id"`
	BusinessID int       `json:"businessId"`
	Key        string    `json:"key"`
	Value      string    `json:"value"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
