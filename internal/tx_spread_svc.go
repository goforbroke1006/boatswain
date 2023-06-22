package internal

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/goforbroke1006/boatswain/domain"
)

func NewTransactionSpreadInfoService(p2pPubSub *pubsub.PubSub) domain.TransactionSpreadInfoService {
	return &txSpreadSvc{p2pPubSub: p2pPubSub}
}

var _ domain.TransactionSpreadInfoService = (*txSpreadSvc)(nil)

type txSpreadSvc struct {
	p2pPubSub *pubsub.PubSub
}

func (svc txSpreadSvc) Spread(tx domain.Transaction) error {
	//TODO implement me
	panic("implement me")
}
