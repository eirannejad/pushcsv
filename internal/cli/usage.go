package cli

import (
	"fmt"
	"os"
	"strings"
)

const version string = "1.0-beta"
const help string = `Push CSV data to database.

Usage:
	pushcsv <db_uri> <table> <file> [--headers] [--purge] [--map=<field_maps>]... [--debug] [--trace] [--dry-run]
	pushcsv <db_uri> <table> --purge

Options:
	-h --help            Show this screen.
	--version            Show version.
	--headers            Indicate first line of csv file is column headers
	--purge              Purge all exiting data from table or collection first
	--map=<field_maps>   Map a source column to table field [from:to]
	                     if mapping is used, all columns missing a map will be ignored.
	--debug              Print debug info
	--trace              Print trace info e.g. full sql queries
	--dry-run            Do everything except commiting data to db

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
