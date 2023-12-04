package models

import "time"

type ServiceType struct {
	ID          int       `json:"id"`
	BusinessID  int       `json:"businessId"`
	Name        string    `json:"name"`
	Icon        string    `json:"icon"`
	Description string    `json:"description"`
	Position    int       `json:"position"`
	CreatedAt   time.Time `json:"createdAt"`
}
