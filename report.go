package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"google.golang.org/api/sheets/v4"
)

type reportBase struct {
	spreadsheetId string
	client        *http.Client
}

type LaneReport struct {
	reportBase
	hoursByTag map[string]float64
	tabId      string
}

func NewLaneReport(spreadsheetId string, client *http.Client, hoursByTag map[string]float64, tabId string) LaneReport {
	return LaneReport{
		reportBase: reportBase{
			spreadsheetId: spreadsheetId,
			client:        client,
		},
		hoursByTag: hoursByTag,
		tabId:      tabId,
	}
}

func (r LaneReport) Update() error {
	srv, err := sheets.New(r.client)
	if err != nil {
		return err
	}

	// Location of the tag cells
	rowOffset := 4
	noTags := 8
	cellRange := fmt.Sprintf("A%v:A%v", rowOffset, rowOffset+noTags-1)
	readRange := fmt.Sprintf("%s!%s", r.tabId, cellRange)

	resp, err := srv.Spreadsheets.Values.Get(r.spreadsheetId, readRange).Do()
	if err != nil {
		return err
	}

	if len(resp.Values) == 0 {
		return fmt.Errorf("No data found in sheet.")
	}

	for idx, row := range resp.Values {
		tag, ok := row[0].(string)
		if !ok {
			return fmt.Errorf("Tag must be of type string.")
		}

		hours, ok := r.hoursByTag[tag]
		if !ok {
			hours = 0.0
		}

		// One cell right we write the values for the cells
		var vr sheets.ValueRange
		myval := []interface{}{hours}
		vr.Values = append(vr.Values, myval)
		writeRange := fmt.Sprintf("%s!B%d", r.tabId, idx+rowOffset)
		_, err := srv.Spreadsheets.Values.Update(r.spreadsheetId, writeRange, &vr).ValueInputOption("RAW").Do()
		if err != nil {
			return err
		}

		fmt.Printf("%v: %v %v\n", idx, tag, hours)
	}

	return nil
}

type HoursReport struct {
	reportBase
	entries []hourTagEntry
	tabId   string
}

func NewHoursReport(spreadsheetId string, client *http.Client, entries []hourTagEntry, tabId string) HoursReport {
	return HoursReport{
		reportBase: reportBase{
			spreadsheetId: spreadsheetId,
			client:        client,
		},
		entries: entries,
		tabId:   tabId,
	}
}

func (r HoursReport) Update() error {
	srv, err := sheets.New(r.client)
	if err != nil {
		return err
	}

	sort.Sort(sort.Reverse(hoursSortedEntries(r.entries)))

	maxEntries := 50

	// Location of the cells
	rowOffset := 18

	entriesLen := len(r.entries)
	for idx := 0; idx < maxEntries; idx++ {
		tag := ""
		hours := ""
		if idx < entriesLen {
			hours = strconv.FormatFloat(r.entries[idx].Hours, 'f', 2, 64)
			tag = r.entries[idx].Tag
		}

		mytag := []interface{}{tag}
		myval := []interface{}{hours}
		fmt.Printf("%v: %v %v\n", idx, tag, hours)

		var tr sheets.ValueRange
		tr.Values = append(tr.Values, mytag)
		writeRange := fmt.Sprintf("%s!G%d", r.tabId, idx+rowOffset)
		_, err := srv.Spreadsheets.Values.Update(r.spreadsheetId, writeRange, &tr).ValueInputOption("RAW").Do()
		if err != nil {
			return err
		}

		// One cell right we write the values for the cells
		var vr sheets.ValueRange
		vr.Values = append(vr.Values, myval)
		writeRange = fmt.Sprintf("%s!H%d", r.tabId, idx+rowOffset)
		_, err = srv.Spreadsheets.Values.Update(r.spreadsheetId, writeRange, &vr).ValueInputOption("RAW").Do()
		if err != nil {
			return err
		}
	}

	return nil
}

type hoursSortedEntries []hourTagEntry

func (e hoursSortedEntries) Len() int {
	return len(e)
}

func (e hoursSortedEntries) Less(i, j int) bool {
	return e[i].Hours < e[j].Hours
}

func (e hoursSortedEntries) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
