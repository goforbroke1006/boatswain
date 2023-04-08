package domain

import "github.com/google/uuid"

type TransactionPayload struct {
	Blockchain    string    `json:"blockchain"`
	ID            uuid.UUID `json:"id"`
	PeerSender    string    `json:"peer_sender"`
	PeerRecipient string    `json:"peer_recipient"`
	Content       string    `json:"content"`
	Timestamp     int64     `json:"timestamp"`
}
