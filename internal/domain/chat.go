package domain

import "time"

type Chat struct {
	ID         int       `json:"id"`
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id,omitempty"`
	GroupID    int       `json:"group_id,omitempty"`
	Message    string    `json:"message"`
	Timestamp  time.Time `json:"timestamp"`
}