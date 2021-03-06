package main

import (
	"encoding/csv"
	"os"
	"strconv"
)

func ParseLanesFile(filename string) (map[string]float64, error) {
	hoursByTag := make(map[string]float64)

	csvFile, err := os.Open(filename)
	if err != nil {
		return hoursByTag, err
	}

	r := csv.NewReader(csvFile)
	records, err := r.ReadAll()
	if err != nil {
		return hoursByTag, err
	}

	for _, record := range records[1:] {
		v, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return hoursByTag, err
		}

		hoursByTag[record[0]] = v
	}

	return hoursByTag, nil
}

type hourTagEntry struct {
	Hours float64
	Tag   string
}

func ParseHoursFile(filename string) ([]hourTagEntry, error) {
	ret := make([]hourTagEntry, 0)

	csvFile, err := os.Open(filename)
	if err != nil {
		return ret, err
	}

	r := csv.NewReader(csvFile)
	records, err := r.ReadAll()
	if err != nil {
		return ret, err
	}

	for _, record := range records[1:] {
		v, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return ret, err
		}

		entry := hourTagEntry{
			Tag:   record[0],
			Hours: v,
		}

		ret = append(ret, entry)
	}

	return ret, nil

}
