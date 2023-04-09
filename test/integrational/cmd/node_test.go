package cmd

//go:generate go install github.com/golang/mock/mockgen@v1.6.0
//go:generate mockgen -source=./../../../domain/block.go -package=mocks -destination=mocks/block.mock.go

import "testing"

func TestNodeReconciliation(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integration tests")
	}

	// TODO: start 2 nodes
	// TODO: fill both node's databases with blocks
	// TODO: start 1 node with empty DB
	// TODO: wait and check last node - should have all blocks
}
