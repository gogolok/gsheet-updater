package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	client, err := NewClient()
	if err != nil {
		log.Fatalln(err)
	}

	filename := os.Getenv("FILE")
	if len(filename) < 1 {
		log.Fatalf("Environment variable FILE must be set.")
	}

	hoursByTag, err := ParseHoursByTagFile(filename)
	if err != nil {
		log.Fatalf("Failed to parse CSV file with hours per tag: %v", err)
	}

	tabId := os.Getenv("TAB_ID")
	if len(tabId) < 1 {
		log.Fatalf("Environment variable TAB_ID must be set.")
	}

	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
	// -> spreadsheetId = 1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms
	spreadsheetId := os.Getenv("SPREADSHEET_ID")
	if len(spreadsheetId) < 1 {
		log.Fatalf("SPREADSHEET_ID not set")
	}

	report := NewLaneReport(spreadsheetId, client, hoursByTag, tabId)
	err = report.Update()
	if err != nil {
		log.Fatalf("Unable to update report: %v", err)
	}
}
