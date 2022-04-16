package main

import (
	"flag"
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
)

const VERSION = "0.0.1"

func main() {
	var (
		url         = flag.String("url", "", "URL to fetch")
		accessToken = flag.String("access token", "", "Optional access token")
		version     = flag.Bool("version", false, "Print version and exit")
	)
	flag.Parse()
	if *version {
		fmt.Printf("HealthChecker %s", VERSION)
		return
	}
	if *url == "" {
		fmt.Println("ERROR : URL required")
		os.Exit(5)
	}
	request := resty.New().R().
		SetHeader("user-agent", fmt.Sprintf("HealthChecker %s", VERSION))
	if *accessToken != "" {
		request.SetHeader("Authorization", fmt.Sprintf("Bearer %s", *accessToken))
	}
	resp, err := request.Get(*url)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(10)
	}
	if !resp.IsSuccess() {
		fmt.Printf("Error: %s", resp.Status())
	}
	return
}
