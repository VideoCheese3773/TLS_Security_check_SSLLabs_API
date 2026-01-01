// Eduardo Mejia
//TODO: Write a program in Go that uses the SSL labs API to check the TLS security of a given domain. Everything is open, so you are free to implement this as you want.

package main

import (
	"fmt" // format for printing and scan the URL
	//"io"       // To read the returned info
	"encoding/json" // to decode the API response
	"net/http"      // this is to make the API call
	"os"            // to work as breaks
	"time"          // to wait for the results of the API call
)

// API URLs
const apiUrl = "https://api.ssllabs.com/api/v4/"
const analyzeUrl = "analyze?host="
const newRequestUrl = "&startNew=on&all=done"
const pollingUrl = "&all=done"

// Data structures for the API, maps the json to Go structs so its easier to decode
type Host struct {
	Host      string     `json:"host"`
	Port      int        `json:"port"`
	Protocol  string     `json:"protocol"`
	Status    string     `json:"status"`
	StartTime int        `json:"startTime"`
	TestTime  int        `json:"testTime"`
	Endpoints []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	IPAddress     string `json:"ipAddress"`
	ServerName    string `json:"serverName"`
	StatusMessage string `json:"statusMessage"`
	Grade         string `json:"grade"`
	HasWarnings   bool   `json:"hasWarnings"`
}

// Functions for the program
func TLSSecurityCheck(url string, email string) (*Host, error) { // First requesto to the API
	fmt.Println("URL to analyze:", url)
	fmt.Println("E-mail used:", email)

	//startNew starts an assesment, all=done only returns info when assesment is done
	var startUrl = apiUrl + analyzeUrl + url + newRequestUrl
	fmt.Println("\n URL for API:", startUrl)

	req, err := http.NewRequest("GET", startUrl, nil) //Can't use http.Get because we need to append the registered email
	if err != nil {                                   // This is to make sure there are no errors in getting stuff from the API, nil is the go version of null btw
		return nil, fmt.Errorf("couldn't create request: %v", err) //Returns an error message
	}
	req.Header.Add("email", email) //Appending email in header of request

	//Using client to send the request
	client := &http.Client{Timeout: 120 * time.Second} //Setting timeout to 120 seconds
	resp, err := client.Do(req)                        //Making the API call
	if err != nil {                                    //Again checking for errors
		return nil, fmt.Errorf("couldn't make the request: %v", err)
	}
	defer resp.Body.Close() //Defer waits until the whole func ends so the body gets closed

	if resp.StatusCode != 200 { //Checking if API call was a success
		return nil, fmt.Errorf("couldn't call API, status code: %d", resp.StatusCode)
	}

	var host Host                                                    //makes this var a Host struct
	if err := json.NewDecoder(resp.Body).Decode(&host); err != nil { // This is to decode the API response and put it into the Host struct
		return nil, fmt.Errorf("couldn't decode API response: %v", err)
	}

	fmt.Println("Assesment Started, Host:", host.Host, "Status:", host.Status)

	return pollResult(url, email)
}

func pollResult(url string, email string) (*Host, error) { // Checks if request is done
	var pollUrl = apiUrl + analyzeUrl + url + pollingUrl
	fmt.Println("URL for polling:", pollUrl)
	cycles := 0

	for {
		req, err := http.NewRequest("GET", pollUrl, nil)
		if err != nil {
			return nil, fmt.Errorf("couldn't create poll request: %v", err)
		}
		req.Header.Add("email", email)

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("couldn't make the poll request: %v", err)
		}

		var host Host
		if err := json.NewDecoder(resp.Body).Decode(&host); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("couldn't decode API poll response: %v", err)
		}
		resp.Body.Close() // Basically same thing as above to check the result of the API call

		fmt.Println("\n Request Status:", host.Status)

		if host.Status == "READY" || host.Status == "ERROR" { // Checking the status of the request
			return &host, nil
		}

		cycles++
		fmt.Printf("Current time: %s, Cycle No: %d \n", time.Now().Format("15:04:05"), cycles)

		if len(host.Endpoints) > 0 { // checks if there are endpoints
			for i, endpoint := range host.Endpoints {
				fmt.Printf("Endpoint %d: %s \n", i+1, endpoint.StatusMessage) // shows progress of each endpoint
			}
		}

		if host.Status == "DNS" || host.Status == "IN_PROGRESS" {
			fmt.Println("Waiting 15 seconds before checking again...")
			time.Sleep(15 * time.Second)
		} else {
			fmt.Println("Waiting 30 seconds before checking again...")
			time.Sleep(30 * time.Second)
		}

		if cycles > 60 { // so that it doesn't run forever
			return nil, fmt.Errorf("Timeout after 60 cycles")
		}

	}
}

func main() {
	var url string // saving url to check
	fmt.Println("\n Please write the URL to check their TLS Security")
	fmt.Scan(&url) // This is how you save inputs from the console to a variable

	var email string // saving email registered in API
	fmt.Println("\n Please write a registered e-mail in the API")
	fmt.Scan(&email)

	host, err := TLSSecurityCheck(url, email)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("\n Test results")
	fmt.Println("Host:", host.Host)
	fmt.Println("Status:", host.Status)

	if len(host.Endpoints) > 0 {
		fmt.Println("Endpoints:")
		for i, endpoint := range host.Endpoints {
			fmt.Printf("\n Endpoint %d \n", i+1)
			fmt.Printf("IP Address: %s \n", endpoint.IPAddress)
			fmt.Printf("Server Name: %s \n", endpoint.ServerName)
			fmt.Printf("TLS Security Grade: %s \n", endpoint.Grade)
			fmt.Printf("Status: %s \n", endpoint.StatusMessage)
			fmt.Printf("Has Warnings: %v \n", endpoint.HasWarnings)
		}
	} else {
		fmt.Println("No endpoints found")
	}
}
