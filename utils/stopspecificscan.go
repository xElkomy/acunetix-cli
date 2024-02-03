package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func stopSpecificScan(target string) {
	url := fmt.Sprintf("%s/api/v1/scans?q=status:processing;", tarURL)
	resp, err := makeRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("[!] Error getting active scans: %v\n", err)
		os.Exit(1)
	}

	type ScanInfo struct {
		ScanID string                       `json:"scan_id"`
		Target struct{ Description string } `json:"target"`
	}

	var scanList []ScanInfo
	if err := json.Unmarshal(resp, &scanList); err != nil {
		fmt.Printf("[!] Error parsing active scans: %v\n", err)
		os.Exit(1)
	}

	for _, scan := range scanList {
		if target == scan.Target.Description {
			if err := stopScan(scan.ScanID); err != nil {
				fmt.Printf("[!] Error stopping scan: %v\n", err)
			}
		}
	}
}
