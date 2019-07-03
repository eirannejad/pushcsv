package cli

import (
	"fmt"
	"os"
)

const version string = "1.6"
const help string = `Push csv/tsv data to database

Usage:
	pushcsv <db_uri> <table> <file> [--headers] [--purge] [--compact] [--map=<field_maps>]... [--debug] [--trace] [--dry-run]
	pushcsv <db_uri> <table> --purge [--debug] [--trace] [--dry-run]

Options:
	-h --help            show this screen
	-V --version         show version
	--headers            when first line of csv file is column headers
	--purge              purge all exiting data from table or collection before pushing new data
	--compact            compact after pushing data (sql vacuum, mongodb compact)
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
	sqlserver:           using github.com/denisenkom/go-mssqldb
	sqlite3:             using github.com/mattn/go-sqlite3

Examples:
	pushcsv postgres://user:pass@data.mycompany.com/mydb users ~/users.csv --headers --purge
	pushcsv mongodb://user:pass@localhost:27017/mydb users ~/users.csv --map=name:fullname --map=email:userid
	pushcsv "mysql:user:pass@tcp(localhost:3306)/tests" users ~/users.csv --purge --map=name:fullname --map=email:userid
	pushcsv sqlserver://user:pass@my-azure-db.database.windows.net?database=mydb users ~/users.csv --purge --map=name:fullname --map=email:userid
	pushcsv sqlite3:data.db users ~/users.csv
`

var printHelpAndExit = func(err error, docoptMessage string) {
	if err != nil {
		// if err occured print full help
		// docopt only includes usage section in its message
		fmt.Fprintln(os.Stderr, help)
		os.Exit(1)
	} else {
		// otherwise print whatever docopt says
		// e.g. reporting version
		fmt.Println(docoptMessage)
		os.Exit(0)
	}
}
