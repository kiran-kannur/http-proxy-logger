package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"sync/atomic"
	"time"
)

// Request counter
var reqCounter int32

type DebugTransport struct{}

func (DebugTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	counter := atomic.AddInt32(&reqCounter, 1)

	startTime := time.Now() // Record the start time

	requestDump, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		return nil, err
	}
	log.Printf("---REQUEST %d---\n\n%s\n\n", counter, string(requestDump))

	response, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close() // Make sure to close the response body

	responseDump, err := httputil.DumpResponse(response, true)
	if err != nil {
		// copying the response body did not work
		return nil, err
	}

	elapsedTime := time.Since(startTime).Milliseconds() // Calculate the elapsed time

	// Read and log the response body
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("---RESPONSE %d--- (Time: %d ms)\n\n%s\n\nResponse Payload: %s\n\n", counter, elapsedTime, string(responseDump), string(responseBody))
	return response, err
}

