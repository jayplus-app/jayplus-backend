package models

type Payment struct {
	ID          int     `json:"id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}
