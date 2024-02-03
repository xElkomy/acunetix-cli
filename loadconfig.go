package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
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
	tarURL = fmt.Sprintf("%s:%d", url, port)
	headers = map[string]string{
		"X-Auth":       apiKey,
		"Content-Type": "application/json",
	}
}
