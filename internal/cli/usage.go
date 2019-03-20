package cli

import (
	"fmt"
	"os"
	"strings"
)

const version string = "1.0-beta"
const help string = `Push CSV data to a database.

Usage:
	pushcsv <db_uri> <table> <file> [--headers] [--purge] [--map=<field_maps>]... [--debug] [--trace] [--dry-run]
	pushcsv <db_uri> <table> --purge

Options:
	-h --help            Show this screen.
	--version            Show version.
	--headers            Indicate first line of csv file is column headers
	--purge              Purge all data from table or collection first
	--map=<field_maps>   Map a single csv column name to table field [from:to]
	--debug              Print debug info
	--trace              Print trace info
	--dry-run            Do everything except for writing data to database

Examples:
	pushcsv sqllite://data.db users ~/users.csv
	pushcsv postgres://localhost:5432/mydb users ~/users.csv --headers --purge
	pushcsv mongodb://localhost:27017/mydb users ~/users.csv --map=name:fullname --map=email:userid
`

var printHelpAndExit = func(err error, usage string) {
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
