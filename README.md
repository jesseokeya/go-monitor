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
	"time"

	pkg "github.com/jesseokeya/go-monitor/monitor"
)

func main() {
	monitorHealth := pkg.Monitor{
		Email: "jesseokeya@gmail.com",
		Links: map[string]string{
			"pocketloan-api":    "https://pocketloan-api.herokuapp.com/health",
			"pocketloan-auth":   "https://pocketloan-auth.herokuapp.com/health",
			"pocketloan-notify": "https://pocketloan-notify.herokuapp.com/health",
		},
	}
	// Perform Health Checks Every Hour
	monitorHealth.Every(1*time.Second, func() {
		// Does Health Checks 20 times
		monitorHealth.Do(20, func() {
			response := monitorHealth.Request("GET", "")
			expected := monitorHealth.Verify(
				response, `{"message":"application is healthy","status":200}`)
			if expected == false {
				message := "Health Checks Failed, Check Application Logs For More Details"
				// Sends an email to the user letting them know thier service is down
				monitorHealth.Alert(message)
			} else {
				fmt.Println("Endpoints are behaving as expected")
			}
		})
	})
}
```

## Terminal Snippet
![alt text](https://github.com/jesseokeya/go-monitor/blob/master/images/Screen%20Shot%202018-07-07%20at%205.57.25%20AM.png?raw=true)

## Sample Alert Email
![alt text](https://github.com/jesseokeya/go-monitor/blob/master/images/Screen%20Shot%202018-07-07%20at%206.02.59%20AM.png?raw=true)

