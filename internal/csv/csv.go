package csv

import (
	"encoding/csv"
	"github.com/pkg/errors"
	"os"
)

type CsvData struct {
	Records [][]string
	Headers []string
}

func ReadCSV(csvfile string, hasHeaders bool) (*CsvData, error) {
	// open csv file
	csvfilehndlr, err := os.Open(csvfile)
	if err != nil {
		return nil, err
	}

	csvreader := csv.NewReader(csvfilehndlr)
	// csvreader.LazyQuotes = true
	records, err := csvreader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) < 1 {
		return nil, errors.New("csv is empty")
	}
	if hasHeaders {
		return &CsvData{Records: records[1:], Headers: records[0]}, nil
	} else {
		return &CsvData{Records: records[1:]}, nil
	}
}
