// Copyrights 2018 Saurabh Badhwar
// The use of this package is goverened by MIT License
// which can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"time"
	"fmt"
	"flag"
	"net/http"
)

var (
	httpClient	*http.Client
)

func init() {
	// Create a new http client that we can use to make requests to
	// our inventory server
	httpClient = &http.Client{
		Timeout: time.Second * 10,
	}
}

// argProcessor initializes the argument processor to take up the arguments 
// that are being provided as an input to the program
func argProcessor() {
	// currently we only listen to the --list argument and don't do
	// any custom processing using that argument.
	flag.Parse()
}

func listProcessor() {
	// Make the request to the endpoint to gather the data from inventory
	resp, err := httpClient.Get("http://localhost:8250/get/inventory")
	if err != nil {
		// we had an error, send it back to the client
		fmt.Fprintf(os.Stdout, "%s", err)
		return
	}
	// Read the response and send it back to the output stream
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(os.Stdout, "%s", data)
}

func main() {
	argProcessor()
}