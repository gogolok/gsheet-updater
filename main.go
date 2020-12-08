package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the root Cobra command
var rootCmd = &cobra.Command{
	Use:   "gsheet-updater",
	Short: "gsheet-udpater is a CLI to update lanes in google docs.",
	Long:  `gsheet-udpater is a CLI to update lanes in google docs.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(newLaneReport())
	rootCmd.AddCommand(newHoursReport())
}

func newLaneReport() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lane",
		Short: "Lane time entries by given tags",
		Long:  `Lane time entries by given tags.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return laneReport()
		},
	}

	return cmd
}

func newHoursReport() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hours",
		Short: "Spent hours by given tags",
		Long:  `Spent hours by given tags.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return hoursReport()
		},
	}

	return cmd
}

func laneReport() error {
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
	return report.Update()
}

func hoursReport() error {
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
