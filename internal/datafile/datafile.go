package datafile

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

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

// https://gist.github.com/bradleypeabody/185b1d7ed6c0c2ab6cec
func DecodeUTF16(b []byte, bigEndian bool) ([]byte, error) {
	if len(b)%2 != 0 {
		return nil, fmt.Errorf("Must have even length byte slice")
	}

	u16s := make([]uint16, 1)

	ret := &bytes.Buffer{}

	b8buf := make([]byte, 4)

	lb := len(b)
	for i := 0; i < lb; i += 2 {
		if bigEndian {
			u16s[0] = uint16(b[i+1]) + (uint16(b[i]) << 8)
		} else {
			u16s[0] = uint16(b[i]) + (uint16(b[i+1]) << 8)
		}
		r := utf16.Decode(u16s)
		n := utf8.EncodeRune(b8buf, r[0])
		ret.Write(b8buf[:n])
	}

	return ret.Bytes(), nil
}

func prepareUTF8(file *os.File, logger *cli.Logger) ([]byte, error) {
	BOM_UTF16_BE := []byte{0xfe, 0xff}
	BOM_UTF16 := []byte{0xff, 0xfe}
	BOM_UTF8 := []byte{0xef, 0xbb, 0xbf}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if bytes.Equal(b[0:3], BOM_UTF8) {
		logger.Debug("Input file encoding detected as UTF8 with BOM")
		return b[3:], nil
	} else if bytes.Equal(b[0:2], BOM_UTF16) {
		logger.Debug("Input file encoding detected as UTF16 (LittleEndian)")
		return DecodeUTF16(b[2:], false)
	} else if bytes.Equal(b[0:2], BOM_UTF16_BE) {
		logger.Debug("Input file encoding detected as UTF16 BigEndian")
		return DecodeUTF16(b[2:], true)
	} else {
		logger.Debug("Input file encoding detected as UTF8")
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

	bytearray, err := prepareUTF8(datafilehndlr, logger)
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
