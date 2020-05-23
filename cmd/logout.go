package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// LogoutCmd logs the user out by cleaning the local state so the user needs to login again.
var LogoutCmd = &cobra.Command{
	Use:   "logout [no options!]",
	Short: "log out the user by removing the user's session from local state.",
	Run: func(cmd *cobra.Command, args []string) {
		Remove("/tmp/getgif.json")

		mag := color.New(color.FgMagenta)
		mag.Printf("You've been logged out!\n")
	},
}
