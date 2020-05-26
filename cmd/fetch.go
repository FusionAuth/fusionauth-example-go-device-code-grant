package cmd

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/fatih/color"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

const (
	giphyApiKey string = "UzwiGmfYRsfoQcDgSS6gz6XtVhRX3siT"
)

type search struct {
	Data       []map[string]interface{} `json:"data"`
	Pagination interface{}              `json:"pagination"`
	Meta       interface{}              `json:"meta"`
}

// FetchCmd provides the Cobra sub command for fetching a gif. Requires the user to be logged in via `login`.
var FetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch a random gif from giphy. User must be logged in before using fetch.",
	Run: func(cmd *cobra.Command, args []string) {
		var user fusionauth.User = fusionauth.User{}
		err := Load("/tmp/getgif.json", &user)
		if err != nil && !os.IsNotExist(err) {
			log.Fatal(err)
		}

		if user.Id == "" {
			red := color.New(color.FgRed, color.Bold)
			red.Printf("You're not logged in yet! Please use `getgif login` before using `fetch`.\n")
		} else {
			gifURL, err := fetchGif()
			if err != nil {
				log.Fatal(err)
			}

			open.Run(gifURL)
		}
	},
}

func fetchGif() (gifURL string, err error) {
	giphyURL, err := url.Parse("https://api.giphy.com/v1/gifs/search")
	if err != nil {
		return "", err
	}

	// Build our Giphy query
	params := url.Values{}
	params.Add("q", "gopher")
	params.Add("api_key", giphyApiKey)
	params.Add("limit", "10")
	giphyURL.RawQuery = params.Encode()

	resp, err := http.Get(giphyURL.String())
	if err != nil {
		return "", err
	}

	var data *search = &search{}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(data)

	// Randomize the gif that we find.
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(10)

	// Walk the JSON response down to the URL of the found gif.
	gif := data.Data[idx]
	gifImg := gif["images"].(map[string]interface{})
	originalImg := gifImg["original"].(map[string]interface{})
	gifURL = originalImg["url"].(string)

	return gifURL, nil
}
