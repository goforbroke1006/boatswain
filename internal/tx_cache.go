package internal

import "github.com/goforbroke1006/boatswain/domain"

func NewTransactionCache() domain.TransactionCache {
	return &txCache{}
}

var _ domain.TransactionCache = (*txCache)(nil)

type txCache struct {
}

func (c txCache) Append(tx domain.Transaction) {
	//TODO implement me
	panic("implement me")
}

func (c txCache) GetFirstN(count uint) ([]domain.Transaction, error) {
	//TODO implement me
	panic("implement me")
}
