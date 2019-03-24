package cli

import (
	"fmt"
	"os"
	"strings"
)

const version string = "1.0"
const help string = `Push csv/tsv data to database

Usage:
	pushcsv <db_uri> <table> <file> [--headers] [--purge] [--map=<field_maps>]... [--debug] [--trace] [--dry-run]
	pushcsv <db_uri> <table> --purge [--debug] [--trace] [--dry-run]

Options:
	-h --help            show this screen
	--version            show version
	--headers            when first line of csv file is column headers
	--purge              purge all exiting data from table or collection before pushing new data
	--map=<field_maps>   map a source column header to table field [from:to]
	                     if mapping is used, all columns missing a map will be ignored,
	                     using mapping is a great way to selectively push data
	                     --map assumes --headers. no need to specify both
	--debug              print debug info
	--trace              print trace info e.g. full sql queries
	--dry-run            do everything except commiting data to db

Supports:
	postgresql:          using github.com/lib/pq
	mongodb:             using gopkg.in/mgo.v2
	mysql:               using github.com/go-sql-driver/mysql
	sqlite3:             using github.com/mattn/go-sqlite3

Examples:
	pushcsv postgres://user:pass@data.mycompany.com/mydb users ~/users.csv --headers --purge
	pushcsv mongodb://localhost:27017/mydb users ~/users.csv --map=name:fullname --map=email:userid
	pushcsv "mysql:ein:test@tcp(localhost:3306)/tests" users ~/users.csv --purge --map=name:fullname --map=email:userid
	pushcsv sqlite3:data.db users ~/users.csv
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
