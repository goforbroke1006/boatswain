package domain

import "github.com/google/uuid"

type Transaction struct {
	metaSenderPeerID string

	ID            uuid.UUID `json:"id"`
	PeerSender    string    `json:"peer_sender"`
	PeerRecipient string    `json:"peer_recipient"`
	Content       string    `json:"content"`
	Timestamp     int64     `json:"timestamp"`
}

func (tp *Transaction) SetSender(peerID string) {
	tp.metaSenderPeerID = peerID
}

func (tp *Transaction) GetSender() string {
	return tp.metaSenderPeerID
}
