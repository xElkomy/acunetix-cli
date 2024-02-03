package main

import (
	"bufio"
	"os"
	"strings"
)

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
