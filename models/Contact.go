package models

type Contact struct {
	ID              int    `json:"id"`
	ContactableType string `json:"contactableType"`
	ContactableID   int    `json:"contactableId"`
	Address         string `json:"address"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
}
