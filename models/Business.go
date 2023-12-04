package models

import "time"

type Business struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	BusinessName string    `json:"businessName"`
	Timezone     string    `json:"timezone"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
