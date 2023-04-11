package main

import (
	"go.uber.org/zap"

	"github.com/goforbroke1006/boatswain/cmd"
)

func main() {
	logger, _ := zap.NewDevelopment() // TODO: select on build
	defer func() { _ = logger.Sync() }()
	zap.ReplaceGlobals(logger)

	cmd.Execute()
}
