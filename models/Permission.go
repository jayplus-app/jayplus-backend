package models

type Permission struct {
	ID     int    `json:"id"`
	Action string `json:"action"`
}
