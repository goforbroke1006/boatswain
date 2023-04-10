package messaging

type Income interface {
	SetSender(peerID string)
	GetSender() string
}
