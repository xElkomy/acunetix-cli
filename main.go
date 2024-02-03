package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/asaskevich/govalidator"
	// Updated import path
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
		"full":  "11111111-1111-1111-1111-111111111111",
		"high":  "11111111-1111-1111-1111-111111111112",
		"weak":  "11111111-1111-1111-1111-111111111115",
		"crawl": "11111111-1111-1111-1111-111111111117",
		"xss":   "11111111-1111-1111-1111-111111111116",
		"sql":   "11111111-1111-1111-1111-111111111113",
	}
)

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

func stopScan(scanID string) error {
	url := fmt.Sprintf("%s/api/v1/scans/%s/abort", tarURL, scanID)
	_, err := makeRequest("POST", url, nil)
	if err != nil {
		return err
	}
	fmt.Printf("[-] Scan stopped, Scan ID: %s\n", scanID)
	return nil
}
