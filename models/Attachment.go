package models

type Attachment struct {
	ID              int    `json:"id"`
	CommunicationID int    `json:"communicationId"`
	URL             string `json:"url"`
}
