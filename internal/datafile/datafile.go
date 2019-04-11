package datafile

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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

func prepareUTF8(file *os.File) ([]byte, error) {
	utf8BOM := []byte{0xef, 0xbb, 0xbf}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if bytes.Equal(b[0:3], utf8BOM) {
		return b[3:], nil
	} else {
		return b, nil
	}
}

func ReadData(dataFile string, options *cli.Options, logger *cli.Logger) (*TableData, error) {
	// open data file
	datafilehndlr, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}
	defer datafilehndlr.Close()

	bytearray, err := prepareUTF8(datafilehndlr)
	if err != nil {
		return nil, err
	}

	// pick/configure reader by file extension
	fileExt := strings.ToLower(filepath.Ext(dataFile))
	var records [][]string
	var readErr error
	if fileExt == ".csv" || fileExt == ".tsv" {
		csvreader := csv.NewReader(strings.NewReader(string(bytearray)))
		// csvreader.LazyQuotes = true
		if fileExt == ".tsv" {
			csvreader.Comma = '\t'
		}
		records, readErr = csvreader.ReadAll()
		if readErr != nil {
			return nil, readErr
		}
	} else {
		return nil, errors.New(fmt.Sprintf("%s file format not supported", fileExt))
	}

	// check if there is data
	if len(records) < 1 {
		return nil, errors.New("data file is empty")
	}

	// if headers exist
	if options.HasHeaders {
		// if maps exits process maps
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
		} else
		// otherwise return all
		{
			return &TableData{
				Name:    options.Table,
				Headers: records[0],
				Records: records[1:],
			}, nil
		}
	} else
	// or just return the dumb data
	{
		return &TableData{
			Name:    options.Table,
			Records: records,
		}, nil
	}
}
