// Eduardo Mejia
//TODO: Write a program in Go that uses the SSL labs API to check the TLS security of a given domain. Everything is open, so you are free to implement this as you want.

package main

import (
	"fmt" // format for printing and scan the URL
	//"io"       // To read the returned info
	//"net/http" // this is to make the API call
	//"os"
)

// API URLs
const apiUrl = "https://api.ssllabs.com/api/v4/"
const analyzeUrl = "analyze?host="

func TLSSecurityCheck(url string, email string) {
	fmt.Println("URL to analyze:", url)
	fmt.Println("E-mail used:", email)

	//startNew starts an assesment, all=done only returns info when assesment is done
	var startUrl = "" + apiUrl + analyzeUrl + url + "&startNew=on&all=done"
	fmt.Println(startUrl)

	/*resp, err := http.Get(startUrl)
	if err != nil { // This is to make sure there are no errors in getting stuff from the API, nil is the go version of null btw
		fmt.Println(err.Error()) //If there is an error, this will print it
		os.Exit(1)               //A number means there was an error, 0 is success
	}
	respData, err := io.ReadAll(resp.Body)
	if err != nil { // Same thing as the one above, just making sure there are no errors in the information recived this time
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(string(respData)) //needs string to be readable, otherwise its just numbers
	*/
}

func main() {
	var url string // saving url to check
	fmt.Println("Please write the URL to check their TLS Security")
	fmt.Scan(&url) // This is how you save inputs from the console to a variable

	var email string // saving email registered in API
	fmt.Println("Please write a registered e-mail in the API")
	fmt.Scan(&email)

	if url == "a" {
		url = "https://pokeapi.co/api/v2/pokemon/ditto" // This whole if section is just a placeholder url to test if the API call works
	}

	TLSSecurityCheck(url, email)
}
