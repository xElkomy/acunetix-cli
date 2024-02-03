package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

func createScan(targetURL string, scanType string) error {
	profileID := scanTypes[scanType]
	if profileID == "" {
		profileID = scanTypes["full"]
	}

	data := map[string]interface{}{
		"address":     targetURL,
		"description": targetURL,
		"criticality": "10",
	}

	resp, err := makeRequest("POST", targetURL+"/api/v1/targets", data)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("error parsing target response: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("error parsing target response: %v", err)
	}

	targetID, ok := result["target_id"].(string)
	if !ok {
		return errors.New("invalid target ID in response")
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

	resp, err = makeRequest("POST", targetURL+"/api/v1/scans", scanData)
	if err != nil {
		return fmt.Errorf("error creating scan: %v", err)
	}

	fmt.Printf("Scan response: %s\n", string(resp))
	return nil
}
