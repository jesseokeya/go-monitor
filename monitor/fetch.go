package monitor

import (
	"bytes"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	err     = color.New(color.FgRed, color.Bold)
	info    = color.New(color.FgCyan, color.Bold)
	loading = color.New(color.FgWhite, color.Bold)
	print   = color.New(color.FgMagenta, color.Bold)
)

// Request the response details of url that was fetched or is about to be fetched
type Request struct {
	URL      string
	Response *http.Response
	Error    error
}

// Send sends a valid http request to multiple urls concurently using go routines
func Send(method string, urls []string, body string) []*Request {
	method = strings.ToUpper(method)
	VerifyMethod(method)
	ch := make(chan *Request, len(urls))
	responses := []*Request{}
	if len(urls) > 0 {
		for _, url := range urls {
			go func(url string) {
				info.Printf("Fetching %s \n", url)
				resp, err := HandleRequest(method, url, body)
				ch <- &Request{url, resp, err}
			}(url)
		}

		for {
			select {
			case r := <-ch:
				print.Printf("%s was fetched\n", r.URL)
				responses = append(responses, r)
				if len(responses) == len(urls) {
					return responses
				}
			case <-time.After(50 * time.Millisecond):
				loading.Printf(".")
			}
		}

	}
	return responses
}

// HandleRequest converts a *http.Request to *http.Response
func HandleRequest(method string, url string, body string) (*http.Response, error) {
	if method == "GET" {
		return http.Get(url)
	}
	jsonStr := []byte(body)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	return resp, nil
}

// VerifyMethod checks if the httpmethod passed in is valid for request that is about to be sent out
func VerifyMethod(method string) {
	HTTPMethods := []string{
		"GET",
		"HEAD",
		"POST",
		"PUT",
		"PATCH",
		"DELETE",
		"CONNECT",
		"OPTIONS",
		"TRACE",
	}
	isValid := false
	for _, valid := range HTTPMethods {
		if method == valid {
			isValid = true
		}
	}
	if !isValid {
		err.Printf("Http Request Method: '%s' is not valid \n", method)
		os.Exit(1)
	}
}
