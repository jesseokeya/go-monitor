# go-monitor
A simple http monitoring library that externally and concurrently monitors rest api endpoints periodically 

# Usage

## Install
```
go get -u github.com/jesseokeya/go-monitor/monitor
```
## Example
```go
package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	pkg "github.com/jesseokeya/go-monitor/monitor"
)

func main() {
	// Perform Health Checks Every Hour
	monitorHealth := pkg.Monitor{
		Email: "jesseokeya@gmail.com",
		Links: map[string]string{
			"pocketloan-api":    "https://pocketloan-api.herokuapp.com/health",
			"pocketloan-auth":   "https://pocketloan-auth.herokuapp.com/health",
			"pocketloan-notify": "https://pocketloan-notify.herokuapp.com/health",
		},
	}

	monitorHealth.Every(1*time.Hour, func() {
		monitorHealth.Do(20, func() {
			response := monitorHealth.Request("GET", "")
			expected := monitorHealth.Equal(response, `{"message":"application is healthy","status":200}`)
			if expected == false {
				message := "Health Checks Failed On PocketLoan Monitoring Service"
        // sends an email to the user letting them know thier service is down
				monitorHealth.Alert(message)
			}
		})
	})
}
```
