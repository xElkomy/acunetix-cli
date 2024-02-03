package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func loadConfig() {
	var apiKey, url, port string
	err := godotenv.Load(os.ExpandEnv("$HOME/.acunetixconfig"))
	if err != nil {
		defaultConfig := map[string]string{
			"url":     "https://localhost",
			"port":    "3443",
			"api_key": "API_KEY",
		}
		err = godotenv.Write(defaultConfig, os.ExpandEnv("$HOME/.acunetixconfig"))
		if err != nil {
			fmt.Println("Failed to write default config:", err)
			os.Exit(1)
		}
	}

	apiKey = os.Getenv("api_key")
	url = os.Getenv("url")
	port = os.Getenv("port")

	if apiKey == "" || url == "" || port == "" {
		fmt.Println("Configuration values (api_key, url, or port) are missing or empty.")
		os.Exit(1)
	}

	// Check if the URL has a valid scheme; if not, use 'http://' as the default scheme.
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	tarURL = fmt.Sprintf("%s:%s", url, port)
	headers = map[string]string{
		"X-Auth":       apiKey,
		"Content-Type": "application/json",
	}
}
