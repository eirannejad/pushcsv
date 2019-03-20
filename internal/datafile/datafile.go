package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"github.com/eirannejad/pushcsv/internal/cli"
	"github.com/pkg/errors"
)

type TableData struct {
	Name    string
	Headers []string
	Records [][]string
}

func (td TableData) HasHeaders() bool {
	return td.Headers != nil
}

func ReadData(dataFile string, options *cli.Options, logger *cli.Logger) (*TableData, error) {
	// open data file
	datafilehndlr, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}

	fileExt := filepath.Ext(dataFile)
	var records [][]string
	var readErr error
	if fileExt == ".csv" {
		csvreader := csv.NewReader(datafilehndlr)
		// csvreader.LazyQuotes = true
		// TODO: csvreader.Comma = '\t'
		records, readErr = csvreader.ReadAll()
		if readErr != nil {
			return nil, readErr
		}
		if len(records) < 1 {
			return nil, errors.New("csv is empty")
		}
	} else {
		return nil, errors.New(fmt.Sprintf("%s file format not supported", fileExt))
	}

	if options.HasHeaders {
		// if maps exits
		if len(options.AttrMaps) > 0 {
			// grab the mapped fields only
			mappedFieldIndices := make([]int, 0)
			mappedFields := make([]string, 0)
			for hidx, header := range records[0] {
				mappedField := options.AttrMaps[header]
				if mappedField != "" {
					mappedFields = append(mappedFields, mappedField)
					mappedFieldIndices = append(mappedFieldIndices, hidx)
				}
			}
			// grab values for mapped fields only
			mappedFieldValues := make([][]string, 0)
			for _, record := range records[1:] {
				recordFieldValues := make([]string, 0)
				for _, fidx := range mappedFieldIndices {
					recordFieldValues = append(recordFieldValues, record[fidx])
				}
				mappedFieldValues = append(mappedFieldValues, recordFieldValues)
			}
			return &TableData{
				Name:    options.Table,
				Headers: mappedFields,
				Records: mappedFieldValues,
			}, nil
		} else {
			// otherwise return all
			return &TableData{
				Name:    options.Table,
				Headers: records[0],
				Records: records[1:],
			}, nil
		}
	} else {
		return &TableData{
			Name:    options.Table,
			Records: records[1:],
		}, nil
	}
}
