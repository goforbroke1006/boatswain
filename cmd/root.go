package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{Use: "boatswain"}
	dappCmd = &cobra.Command{Use: "dapp"}
)

func Execute() {
	rootCmd.AddCommand(dappCmd)

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
