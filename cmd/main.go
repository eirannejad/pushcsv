package datafile

import (
	"fmt"
	"os"

	"github.com/eirannejad/pushcsv/internal/cli"
	"github.com/eirannejad/pushcsv/internal/datafile"
	"github.com/eirannejad/pushcsv/internal/persistance"
)

func Run() {
	// process arguments
	argv := os.Args[1:]
	options := cli.NewOptions(argv)

	// log options if requested
	logger := cli.NewLogger(options)
	logger.Debug(options)

	// check connection string and determine target db driver
	// verify db is supported
	// this step should be before starting io on reading data file
	dbConfig, cErr := persistance.NewDatabaseConfig(options.ConnString)
	if cErr != nil {
		fmt.Fprintln(os.Stderr, cErr.Error())
		os.Exit(1)
	}

	// read datafile
	// prepare the data for writer; fixes the data mappings
	tableData, err := datafile.ReadData(options.CsvFile, options, logger)
	if err != nil {
		// TODO: Create Error And Exit Func?
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// request a writer for db
	writer, nErr := persistance.NewWriter(dbConfig, options, logger)
	if nErr != nil {
		fmt.Fprintln(os.Stderr, nErr.Error())
		os.Exit(1)
	}

	// write to db if not dry run
	if !options.DryRun {
		result, wErr := writer.Write(tableData)
		if wErr != nil {
			fmt.Fprintln(os.Stderr, wErr.Error())
			os.Exit(1)
		}
		logger.Print(fmt.Sprintf("Successfully updated %d records.", result.Count))
	}
}
