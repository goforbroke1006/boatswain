package consensus

import "testing"

//go:generate go install github.com/golang/mock/mockgen@v1.6.0
//go:generate mockgen -source=./../../../../domain/block.go -package=mocks -destination=mocks/block.mock.go

func TestGenerator(t *testing.T) {
	// TODO:
}
