package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Create a reusable HTTP client with an insecure transport (for ignoring SSL certificate validation).
var httpClient = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}

func makeRequest(method, url string, data interface{}) ([]byte, error) {
	// Marshal the request data into JSON.
	requestBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Create an HTTP request with the request method, URL, and request body.
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	// Set headers from the global "headers" map.
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Send the HTTP request using the reusable HTTP client.
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body.
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check the HTTP status code.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP Error: %s", resp.Status)
	}

	return responseBody, nil
}
