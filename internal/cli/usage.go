package cli

import (
	"fmt"
	"os"
	"strings"
)

const version string = "1.0-beta"
const help string = `Push CSV data to a database.

Usage:
	pushcsv <db_uri> <table> <csv_file> [--headers] [--purge] [--map=<field_maps>] [--debug]
	pushcsv <db_uri> <table> --purge

Options:
	-h --help            Show this screen.
	--version            Show version.
	--headers            Indicate first line of csv file is column headers
	--purge              Purge all data from table or collection first
	--map=<field_maps>   maps column index to db field in order [first;second;third]

Examples:
	pushcsv sqllite://data.db users ~/users.csv
	pushcsv postgres://localhost:5432/mydb users ~/users.csv --headers --purge
	pushcsv mongodb://localhost:27017/mydb users ~/users.csv --map=lastname;firstname
`

var PrintHelpAndExit = func(err error, usage string) {
	if err != nil {
		fmt.Fprintln(os.Stderr, help)
		os.Exit(1)
	} else {
		if strings.Index(usage, version) == 0 {
			fmt.Println(usage)
		} else {
			fmt.Println(help)
		}
		os.Exit(0)
	}
}
