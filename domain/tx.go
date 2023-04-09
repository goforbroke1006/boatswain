package domain

import "github.com/google/uuid"

type TransactionPayload struct {
	MetaSenderPeerID string `json:"-"`

	Blockchain    string    `json:"blockchain"`
	ID            uuid.UUID `json:"id"`
	PeerSender    string    `json:"peer_sender"`
	PeerRecipient string    `json:"peer_recipient"`
	Content       string    `json:"content"`
	Timestamp     int64     `json:"timestamp"`
}

func (tp *TransactionPayload) SetSender(peerID string) {
	tp.MetaSenderPeerID = peerID
}

func (tp *TransactionPayload) GetSender() string {
	return tp.MetaSenderPeerID
}
