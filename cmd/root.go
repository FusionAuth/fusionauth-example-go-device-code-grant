/*
Copyright © 2020 Matt Gowie <matt@masterpoint.io>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	// APIKey is the FA API Key created in the FA UI under "Settings > API Keys"
	APIKey string

	// ClientID is the OAuth client_id of our FA Application
	ClientID string

	// TenantID is the ID of our FA instance's default Tenant
	TenantID string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "getgif",
	Short: "A small CLI tool to fetch Gifs only once the user is authenticated.",
	Long: `
A small CLI tool to fetch Gifs only once the user is authenticated.
This is an example application showing off the FusionAuth service and their golang client library.
You can read the accompanying blog post @ https://fusionauth.io/blog/2020/06/17/building-cli-app-with-device-grant-and-golang
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	godotenv.Load(".env")
	ClientID = os.Getenv("FA_CLIENT_ID")
	TenantID = os.Getenv("FA_TENANT_ID")
	APIKey = os.Getenv("FA_API_KEY")

	if ClientID == "TODO" || TenantID == "TODO" || APIKey == "TODO" {
		red := color.New(color.FgRed, color.Bold)
		red.Printf("You need to finish the setup and add your FA configuration to `.env`. Do so and come back!\n")
		os.Exit(1)
	}

	rootCmd.AddCommand(LoginCmd)
	rootCmd.AddCommand(LogoutCmd)
	rootCmd.AddCommand(FetchCmd)
}
