package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

var VERSION = "Unreleased"

type rawHeaders []string

func (r *rawHeaders) String() string {
	return "*Something Important*"
}

func (r *rawHeaders) Set(s string) error {
	*r = append(*r, s)
	return nil
}

func request(method *string, url *string, headers *rawHeaders, timeout *int) (*http.Response, error) {
	if *url == "" {
		fmt.Println("Error : URL required")
		os.Exit(5)
	}
	request, _ := http.NewRequest(*method, *url, nil)
	addHeaders(request, headers)
	return (&http.Client{
		Timeout: time.Second * time.Duration(*timeout),
	}).Do(request)
}

func addHeaders(request *http.Request, headers *rawHeaders) {
	for _, header := range *headers {
		parsedHeader := strings.Split(header, ":")
		request.Header.Set(parsedHeader[0], parsedHeader[1])
	}
}

func read(closer io.ReadCloser) []byte {
	d, _ := io.ReadAll(closer)
	_ = closer.Close()
	return d
}

func main() {
	var (
		headers rawHeaders
		url     = flag.String("url", "", "URL to fetch")
		version = flag.Bool("version", false, "Print version and exit")
		verbose = flag.Bool("verbose", false, "Logging verbosely")
		timeout = flag.Int("timeout", 15, "Timeout")
		method  = flag.String("method", "GET", "HTTP Method")
	)
	flag.Var(&headers, "headers", "headers to send")
	flag.Parse()
	if *version {
		fmt.Printf("HealthChecker %s\n", VERSION)
		return
	}
	resp, err := request(method, url, &headers, timeout)
	if err != nil {
		fmt.Printf("Error : %s\n", err)
		os.Exit(10)
	}
	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		if *verbose == true {
			fmt.Printf("Status : %s\nResponse Body : %s\n", resp.Status, read(resp.Body))
		} else {
			fmt.Printf("Error : %s\n", resp.Status)
		}
		os.Exit(resp.StatusCode)
	} else {
		if *verbose == true {
			fmt.Printf("Status : %s\nResponse Body : %s\n", resp.Status, read(resp.Body))
		}
	}
	return
}
