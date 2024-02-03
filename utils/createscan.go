package main

import (
	"encoding/json"
	"fmt"
)

func createScan(targetURL string, scanType string) {
	profileID = scanTypes[scanType]
	if profileID == "" {
		profileID = scanTypes["full"]
	}

	data := map[string]interface{}{
		"address":     targetURL,
		"description": targetURL,
		"criticality": "10",
	}

	resp, err := makeRequest("POST", tarURL+"/api/v1/targets", data)
	if err != nil {
		fmt.Printf("Error creating target: %v\n", err)
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		fmt.Printf("Error parsing target response: %v\n", err)
		return
	}

	targetID, ok := result["target_id"].(string)
	if !ok {
		fmt.Println("Invalid target ID in response")
		return
	}

	fmt.Printf("[*] Running scan on: %s\n", targetURL)

	scanData := ScanData{
		TargetID:  targetID,
		ProfileID: profileID,
		Schedule: struct {
			Disable       bool   `json:"disable"`
			StartDate     string `json:"start_date,omitempty"`
			TimeSensitive bool   `json:"time_sensitive"`
		}{
			Disable:       false,
			StartDate:     "",
			TimeSensitive: false,
		},
	}

	resp, err = makeRequest("POST", tarURL+"/api/v1/scans", scanData)
	if err != nil {
		fmt.Printf("Error creating scan: %v\n", err)
		return
	}

	fmt.Printf("Scan response: %s\n", string(resp))
}
