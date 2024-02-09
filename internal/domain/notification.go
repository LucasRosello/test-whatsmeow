package domain

import "time"

type Notification struct {
	ID         string    `db:"id"`
	Sender     string    `json:"sender"`
	Receiver   string    `json:"receiver"`
	DateToSend time.Time `json:"date_to_send"`
	DateSended time.Time `json:"date_sended"`
	Message    string    `json:"message"`
	Status     string    `json:"status"`
}
