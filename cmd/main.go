package cmd

import (
	"errors"
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
	logger.Trace(options)
	for key, value := range *options.Opts {
		logger.Debug(fmt.Sprintf("%s=%v", key, value))
	}

	// check connection string and determine target db driver
	// verify db is supported
	// this step should be before starting io on reading data file
	dbConfig, cErr := persistance.NewDatabaseConfig(options.ConnString)
	if cErr != nil {
		errorAndExit(cErr)
	}

	var tableData *datafile.TableData
	var readRrr error
	// check if input file is specified
	if options.DataFile != "" {
		// check if writer needs headers
		if dbConfig.NeedsHeaders && !options.HasHeaders {
			errorAndExit(
				errors.New(
					"headers are required to write to this database. " +
						"make sure source file has headers on first line and " +
						"use --headers flag"))
		}

		// read datafile
		// prepare the data for writer; fixes the data mappings
		tableData, readRrr =
			datafile.ReadData(options.DataFile, options, logger)
		if readRrr != nil {
			errorAndExit(readRrr)
		}
	}

	// request a writer for db
	writer, nErr := persistance.NewWriter(dbConfig, options, logger)
	if nErr != nil {
		errorAndExit(nErr)
	}

	// write to db
	var result *persistance.Result
	var commitErr error
	if tableData != nil {
		result, commitErr = writer.Write(tableData)
	} else {
		result, commitErr = writer.Purge(options.Table)
	}
	if commitErr != nil {
		errorAndExit(commitErr)
	}
	fmt.Println(result.Message)
}

func errorAndExit(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}
