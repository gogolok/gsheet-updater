package main

import (
	"fmt"
	"net/http"

	"google.golang.org/api/sheets/v4"
)

type Report struct {
	spreadsheetId string
	client        *http.Client
	hoursByTag    map[string]float64
	tabId         string
}

func NewReport(spreadsheetId string, client *http.Client, hoursByTag map[string]float64, tabId string) Report {
	return Report{
		spreadsheetId: spreadsheetId,
		client:        client,
		hoursByTag:    hoursByTag,
		tabId:         tabId,
	}
}

func (r Report) Update() error {
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
			return fmt.Errorf("No hours found for tag %v", tag)
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
