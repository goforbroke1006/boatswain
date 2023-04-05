package main

import (
	"go.uber.org/zap"

	"github.com/goforbroke1006/boatswain/cmd"
)

func main() {
	logger, _ := zap.NewProduction()
	defer func() { _ = logger.Sync() }()

	cmd.Execute()
}
