package cmd

import (
	"fmt"
	"github.com/eirannejad/pushcsv/internal/cli"
	"github.com/eirannejad/pushcsv/internal/csv"
	"github.com/eirannejad/pushcsv/internal/persistance"
	"os"
)

func Run() {
	// process arguments
	argv := os.Args[1:]
	options := cli.NewOptions(argv)

	// log options if requested
	logger := cli.NewLogger()
	logger.PrintDebug = options.Debug
	logger.Debug(options)

	// read csv
	csvData, err := csv.ReadCSV(options.CsvFile, options.HasHeaders)
	if err != nil {
		// TODO: Create Error And Exit Func?
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// check connection string and determine target db driver
	writer, nErr := persistance.NewWriter(logger, options)
	if nErr != nil {
		fmt.Fprintln(os.Stderr, nErr.Error())
		os.Exit(1)
	}
	result, wErr := writer.Write(csvData)
	if wErr != nil {
		fmt.Fprintln(os.Stderr, wErr.Error())
		os.Exit(1)
	}

	logger.Print(fmt.Sprintf("Successfully updated %d records.", result.Count))
}
