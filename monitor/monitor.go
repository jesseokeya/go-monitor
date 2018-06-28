package monitor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

// Monitor hold the information required to monitor different urls on interval and
// an email to notify a user on how the application is doing
type Monitor struct {
	Links map[string]string
	Email string
}

// EmailResponse is the stucture of the email response json object
type EmailResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// Request sends a POST, GET, PUT, PATCH, DELETE etc.. requests to all urls specified
func (m *Monitor) Request(method string, data string) []*http.Response {
	return m.Response(method, data)
}

// Urls gets all the urls specified int Monitor.Links
func (m *Monitor) Urls() []string {
	result := []string{}
	for _, value := range m.Links {
		result = append(result, value)
	}
	return result
}

// Keys gets all the urls specified int Monitor.Links
func (m *Monitor) Keys() []string {
	result := []string{}
	for key := range m.Links {
		result = append(result, key)
	}
	return result
}

// Response collects all the http response objects in an array and returns the array
func (m *Monitor) Response(method string, data string) []*http.Response {
	responses := []*http.Response{}
	urls := m.Urls()
	results := Send(method, urls, data)
	for _, result := range results {
		responses = append(responses, result.Response)
	}
	return responses
}

// Do runs any function t amount of times specified
func (m *Monitor) Do(t int, f func()) {
	for i := 0; i < t; i++ {
		f()
	}
}

// Every executes a request to specified urls at intervals you specify. it also
// lets you choose how many times you would like to run that function
func (m *Monitor) Every(d time.Duration, f func()) {
	for range time.Tick(d) {
		f()
	}
}

// Alert sends you an email anytime you application is down or unreacheable
func (m *Monitor) Alert(message string) {
	url := "https://pocketloan-notify.herokuapp.com/v1/api/notify/send-email"
	body := `{"email_message":"` + message + `", "recipientAddress":"` + m.Email + `", "first_name":"pocketloan admin", "template":"tem_rPMJvG9KVwF4XQSwxcgkx3x7"}`
	fmt.Println(body)
	resp, err := HandleRequest("POST", url, body)
	m.handleError(err)

	contents, err := ioutil.ReadAll(resp.Body)
	m.handleError(err)

	expected := EmailResponse{
		Message: "email successfully sent",
		Status:  "200",
	}
	responseObject := EmailResponse{}
	err = json.Unmarshal(contents, &responseObject)
	m.handleError(err)

	if expected != responseObject {
		panic("Something went wrong in trying to send an email")
	}
}

// Verify ensures you response is what is expected
func (m *Monitor) Verify(r []*http.Response, expected string) bool {
	result := false
	for _, resp := range r {
		contents, err := ioutil.ReadAll(resp.Body)
		m.handleError(err)
		result, _ = m.Equal(expected, string(contents))
	}
	return result
}

// Equal compares equality between two strings
func (m *Monitor) Equal(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 1 :: %s", err.Error())
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 2 :: %s", err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}

// handleError panic any error within the monitor package
func (m *Monitor) handleError(err error) {
	if err != nil {
		panic(err)
	}
}
