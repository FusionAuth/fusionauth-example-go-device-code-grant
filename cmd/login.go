package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/fatih/color"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

const (
	host string = "http://localhost:9011"
)

var (
	httpClient = &http.Client{
		Timeout: time.Second * 10,
	}
	baseURL, _ = url.Parse(host)
	faClient   *fusionauth.FusionAuthClient
)

// LoginCmd provides the subcommand for logging into the FA server using the Device Flow.
var LoginCmd = &cobra.Command{
	Use:   "login [no options!]",
	Short: "Login to the FA server using the OAuth Device Flow.",
	Run: func(cmd *cobra.Command, args []string) {
		faClient = fusionauth.NewClient(httpClient, baseURL, APIKey)
		openIDConfig, err := faClient.RetrieveOpenIdConfiguration()
		if err != nil {
			log.Fatal(err)
		}

		deviceResp, err := startDeviceGrantFlow(openIDConfig.DeviceAuthorizationEndpoint)
		if err != nil {
			log.Fatal(err)
		}

		informUserAndOpenBrowser(deviceResp.UserCode)

		accessToken, err := startPolling(openIDConfig.TokenEndpoint, deviceResp.DeviceCode, deviceResp.Interval)
		if err != nil {
			log.Fatal(err)
		}

		fetchAndSaveUser(accessToken)
	},
}

func startDeviceGrantFlow(deviceAuthEndpoint string) (*fusionauth.DeviceResponse, error) {
	var result *fusionauth.DeviceResponse = &fusionauth.DeviceResponse{}

	resp, err := http.PostForm(deviceAuthEndpoint, url.Values{
		"client_id":            {"7dde5f47-5000-4580-8003-b3b8d1cbe2e9"},
		"scope":                {"offline_access"},
		"metaData.device.name": {"Golang CLI App"},
		"metaData.device.type": {string(fusionauth.DeviceType_OTHER)},
	})

	if err != nil {
		return result, err
	}

	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(result)

	return result, nil
}

func informUserAndOpenBrowser(userCode string) {
	cyan := color.New(color.FgCyan)
	cyan.Printf("Your User Code is: ")

	red := color.New(color.FgRed, color.Bold)
	red.Printf("%s\n", userCode)

	fmt.Printf("Opening browser for code entry...\n")

	// Wait a few seconds to give user a chance to check out the printed user code.
	time.Sleep(3 * time.Second)

	url := fmt.Sprintf("%s/oauth2/device?client_id=%s&tenantId=%s", host, ClientID, TenantID)
	open.Run(url)
}

func startPolling(tokenEndpoint string, deviceCode string, retryInterval int) (*fusionauth.AccessToken, error) {
	var result *fusionauth.AccessToken = &fusionauth.AccessToken{}
	yellow := color.New(color.FgYellow, color.Bold)
	blue := color.New(color.FgBlue, color.Bold)

	for {
		resp, err := http.PostForm(tokenEndpoint, url.Values{
			"device_code": {deviceCode},
			"grant_type":  {"urn:ietf:params:oauth:grant-type:device_code"},
			"client_id":   {ClientID},
		})

		if err != nil {
			return result, err
		}

		// 400 status code (StatusBadRequest) is our sign that the user
		// hasn't completed their device login yet, sleep and then continue.
		if resp.StatusCode == http.StatusBadRequest {

			// Sleep for the retry interval and print a dot for each second.
			for i := 0; i < retryInterval; i++ {
				if i == 0 {
					blue.Printf(".")
				} else {
					yellow.Printf(".")
				}
				time.Sleep(time.Second)
			}

			continue
		}

		fmt.Printf("\n")
		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(result)

		return result, nil
	}
}

func fetchAndSaveUser(token *fusionauth.AccessToken) {
	resp, _, err := faClient.RetrieveUserUsingJWT(token.AccessToken)
	if err != nil {
		log.Fatal(err)
	}

	// Save our User object for later usage in fetch
	Save("/tmp/getgif.json", resp.User)

	mag := color.New(color.FgMagenta)
	mag.Printf("You successfully authenticated! You can now use `getgif fetch`!\n")
}
