package models

type SMS struct {
	ID       int    `json:"id"`
	Content  string `json:"content"`
	ToNumber string `json:"toNumber"`
}
