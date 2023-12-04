package models

import "time"

type Communication struct {
	ID        int       `json:"id"`
	ChannelID int       `json:"channelId"`
	UserID    int       `json:"userId"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}
