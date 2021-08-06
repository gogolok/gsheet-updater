package main

import (
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Version = undefinedVersion

const (
	// undefinedVersion should take the form `channel-version`
	undefinedVersion = "dev-undefined"
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
	rootCmd.AddCommand(newCmdVersion())
	rootCmd.AddCommand(newLaneReport())
	rootCmd.AddCommand(newHoursReport())
	rootCmd.AddCommand(newLastRunTimestamp())
}

type versionOptions struct {
	shortVersion bool
}

func newVersionOptions() *versionOptions {
	return &versionOptions{
		shortVersion: false,
	}
}

func newCmdVersion() *cobra.Command {
	options := newVersionOptions()

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the client version information",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			runVersion(options, os.Stdout)
		},
	}

	cmd.PersistentFlags().BoolVar(&options.shortVersion, "short", options.shortVersion, "Print the version number(s) only, with no additional output")

	return cmd
}

func runVersion(options *versionOptions, stdout io.Writer) {
	clientVersion := Version

	if options.shortVersion {
		fmt.Fprintln(stdout, clientVersion)
	} else {
		fmt.Fprintf(stdout, "Client version: %s\n", clientVersion)
	}
}

func newLaneReport() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lane",
		Short: "Spent hours per lane",
		Long:  `Spent hours per lines.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return laneReport()
		},
	}

	return cmd
}

func newHoursReport() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hours",
		Short: "Spent hours per pattern",
		Long:  `Spent hours per pattern.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			maxEntries, _ := cmd.Flags().GetInt("max-entries")
			startColumn, _ := cmd.Flags().GetString("start-column")

			startColumnByte := startColumn[0]

			return hoursReport(maxEntries, startColumnByte)
		},
	}

	cmd.Flags().IntP("max-entries", "m", 50, "Max entries to consider.")
	cmd.Flags().StringP("start-column", "c", "G", "What column to write entries to.")

	return cmd
}

func newLastRunTimestamp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "last-run-timestamp",
		Short: "Write timestamp of last run",
		Long:  `Write the current time as timestamp into the sheet so we're aware when the tool ran the last time`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return lastRunTimestamp()
		},
	}

	cmd.Flags().IntP("max-entries", "m", 50, "Max entries to consider.")
	cmd.Flags().StringP("start-column", "c", "G", "What column to write entries to.")

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

	hoursByTag, err := ParseLanesFile(filename)
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

func hoursReport(maxEntries int, startColumn byte) error {
	client, err := NewClient()
	if err != nil {
		log.Fatalln(err)
	}

	filename := os.Getenv("FILE")
	if len(filename) < 1 {
		log.Fatalf("Environment variable FILE must be set.")
	}

	hoursByTag, err := ParseHoursFile(filename)
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

	report := NewHoursReport(spreadsheetId, client, hoursByTag, tabId, maxEntries, startColumn)
	return report.Update()
}

func lastRunTimestamp() error {
	client, err := NewClient()
	if err != nil {
		log.Fatalln(err)
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

	report := NewLastRunTimestampReport(spreadsheetId, client, tabId)
	return report.Update()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
