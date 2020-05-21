package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// LoginCmd provide the Cobra sub command for logging into the FA server.
var LoginCmd = &cobra.Command{
	Use:   "login [no options!]",
	Short: "Login to the FA server using the OAuth device code grant type.",
	Run: func(cmd *cobra.Command, args []string) {
      fmt.Printf("Inside login - Running with args: %v\n", args)
    },
}