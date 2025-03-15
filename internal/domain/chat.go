package domain

import "time"

type Chat struct {
	ID        string    `json:"id"`
	Sender    string    `json:"sender"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}
