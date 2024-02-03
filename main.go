package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/joho/godotenv"
)

type ScanProfile struct {
	Full  string `json:"full"`
	High  string `json:"high"`
	Weak  string `json:"weak"`
	Crawl string `json:"crawl"`
	XSS   string `json:"xss"`
	SQL   string `json:"sql"`
}

type ScanData struct {
	TargetID  string `json:"target_id"`
	ProfileID string `json:"profile_id"`
	Schedule  struct {
		Disable       bool   `json:"disable"`
		StartDate     string `json:"start_date,omitempty"`
		TimeSensitive bool   `json:"time_sensitive"`
	} `json:"schedule"`
}

var (
	tarURL    string
	headers   map[string]string
	profileID string
	scanTypes = map[string]string{
		"full":    "11111111-1111-1111-1111-111111111111",
		"high":    "11111111-1111-1111-1111-111111111112",
		"weak":    "11111111-1111-1111-1111-111111111115",
		"crawl":   "11111111-1111-1111-1111-111111111117",
		"xss":     "11111111-1111-1111-1111-111111111116",
		"sql":     "11111111-1111-1111-1111-111111111113",
		"xelkomy": os.Getenv("my-scanprofile"),
	}
)

func loadConfig() {
	var apiKey, url, port string
	err := godotenv.Load(os.ExpandEnv("$HOME/.acunetixconfig"))
	if err != nil {
		defaultConfig := map[string]string{
			"my-scanprofile": "YOURSCANPROFILEID",
			"url":            "https://localhost",
			"port":           "3443",
			"api_key":        "API_KEY",
		}
		err = godotenv.Write(defaultConfig, os.ExpandEnv("$HOME/.acunetixconfig"))
	}

	apiKey = os.Getenv("api_key")
	url = os.Getenv("url")
	port = os.Getenv(port)
	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("Failed to parse port number: %v", err)
	}
	tarURL = fmt.Sprintf("%s:%d", url, portInt)
	headers = map[string]string{
		"X-Auth":       apiKey,
		"Content-Type": "application/json",
	}
}
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

func makeRequest(method, url string, data interface{}) ([]byte, error) {
	requestBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP Error: %s", resp.Status)
	}

	return responseBody, nil
}

func main() {
	loadConfig()

	if len(os.Args) < 2 {
		fmt.Println("Usage: acunetix-cli [-h]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "scan":
		if len(os.Args) < 4 {
			fmt.Println("Usage: acunetix-cli scan [-p] [-d DOMAIN | -f FILE] [-t TYPE]")
			os.Exit(1)
		}

		scanType := "full"
		domain := ""
		filePath := ""
		usePipe := false

		for i := 2; i < len(os.Args); i++ {
			arg := os.Args[i]
			switch arg {
			case "-t":
				i++
				scanType = os.Args[i]
			case "-d":
				i++
				domain = os.Args[i]
			case "-f":
				i++
				filePath = os.Args[i]
			case "-p":
				usePipe = true
			}
		}

		if usePipe {
			var urls []string
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				url := scanner.Text()
				if govalidator.IsURL(url) {
					urls = append(urls, url)
				}
			}
			for _, url := range urls {
				createScan(url, scanType)
			}
		} else if domain != "" {
			if govalidator.IsURL(domain) {
				createScan(domain, scanType)
			} else {
				fmt.Println("[!] Invalid URL:", domain)
			}
		} else if filePath != "" {
			targets, err := readTargetsFromFile(filePath)
			if err != nil {
				fmt.Printf("[!] Error reading file: %v\n", err)
				os.Exit(1)
			}
			for _, target := range targets {
				createScan(target, scanType)
			}
		} else {
			fmt.Println("[!] Must provide either domain or file containing list of targets\nFor Help: acunetix-cli scan -h")
		}
	case "stop":
		if len(os.Args) < 3 {
			fmt.Println("Usage: acunetix-cli stop [-d DOMAIN | -a]")
			os.Exit(1)
		}

		domain := ""
		stopAll := false

		for i := 2; i < len(os.Args); i++ {
			arg := os.Args[i]
			switch arg {
			case "-d":
				i++
				domain = os.Args[i]
			case "-a":
				stopAll = true
			}
		}

		if domain != "" {
			stopSpecificScan(domain)
		} else if stopAll {
			stopAllScans()
		} else {
			fmt.Println("[!] Must provide either domain or stop all flag\nFor Help: acunetix-cli stop -h")
		}
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func readTargetsFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var targets []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		target := strings.TrimSpace(scanner.Text())
		if target != "" {
			targets = append(targets, target)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return targets, nil
}

func stopScan(scanID string) error {
	url := fmt.Sprintf("%s/api/v1/scans/%s/abort", tarURL, scanID)
	_, err := makeRequest("POST", url, nil)
	if err != nil {
		return err
	}
	fmt.Printf("[-] Scan stopped, Scan ID: %s\n", scanID)
	return nil
}

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

func stopAllScans() {
	url := fmt.Sprintf("%s/api/v1/scans?q=status:processing;", tarURL)
	resp, err := makeRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("[!] Error getting active scans: %v\n", err)
		os.Exit(1)
	}

	type ScanInfo struct {
		ScanID string `json:"scan_id"`
	}

	var scanList []ScanInfo
	if err := json.Unmarshal(resp, &scanList); err != nil {
		fmt.Printf("[!] Error parsing active scans: %v\n", err)
		os.Exit(1)
	}

	for _, scan := range scanList {
		if err := stopScan(scan.ScanID); err != nil {
			fmt.Printf("[!] Error stopping scan: %v\n", err)
		}
	}
}
