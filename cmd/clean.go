package cmd

import (
	"github.com/spf13/cobra"
)

// CleanCmd cleans the local state so the user needs to login again.
var CleanCmd = &cobra.Command{
	Use:   "clean [no options!]",
	Short: "Clean getgif local state, forcing another login.",
	Run: func(cmd *cobra.Command, args []string) {
		Remove("/tmp/getgif.json")
	},
}
