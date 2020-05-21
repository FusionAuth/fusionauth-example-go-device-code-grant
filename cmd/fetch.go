package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// FetchCmd provides the Cobra sub command for fetching a gif. Requires the user to be logged in via `login`.
var FetchCmd = &cobra.Command{
	Use:   "fetch [no options!]",
	Short: "Fetch a random gif from giphy. User must be logged in before using fetch.",
	Run: func(cmd *cobra.Command, args []string) {
      fmt.Printf("Inside fetch - Running with args: %v\n", args)
    },
}