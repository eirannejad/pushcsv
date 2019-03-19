package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/docopt/docopt-go"
	_ "github.com/lib/pq"
	// _ "github.com/mattn/go-sqlite3"
	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
)

type MyLog struct {
	PrintDebug bool
}

const version string = "1.0"
const postgres_insert string = "INSERT INTO %s values "
const help string = `Push CSV data to a database.

Usage:
	pushcsv <conn> <table> <csv_file> [--skip-first] [--purge] [--map=<field_maps>] [--debug]

Options:
	-h --help            Show this screen.
	--version            Show version.
	--skip-first         Skip first line of CSV file
	--purge              Purge all data from table or collection first
	--map=<field_maps>   maps a csv col name to db col [from:to]

Examples:
	$ pushcsv sqllite://data.db users ~/users.csv
	$ pushcsv postgres://localhost:5432/mydb users ~/users.csv --map=name:fullname
	$ pushcsv mongodb://localhost:27017 users ~/users.csv
`

var ml MyLog = MyLog{}

func (m *MyLog) Debug(args ...interface{}) {
	if m.PrintDebug {
		log.Print(args...)
	}
}

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

func ProcessArgs(opts *docopt.Opts) (string, string, string, bool, bool, bool) {
	connStr, _ := opts.String("<conn>")
	table, _ := opts.String("<table>")
	csvfile, _ := opts.String("<csv_file>")
	skipfirst, _ := opts.Bool("--skip-first")
	purge, _ := opts.Bool("--purge")
	debug, _ := opts.Bool("--debug")
	return connStr, table, csvfile, skipfirst, purge, debug
}

func ReadCSV(csvfile string) ([][]string, error) {
	// open csv file
	csvfilehndlr, err := os.Open(csvfile)
	if err != nil {
		panic(err)
	}

	csvreader := csv.NewReader(csvfilehndlr)
	csvreader.LazyQuotes = true
	return csvreader.ReadAll()
}

func PurgeTable(db *sql.DB, table string) (sql.Result, error) {
	return db.Exec(fmt.Sprintf(`TRUNCATE TABLE %s`, table))
}

func PushCSV(db *sql.DB, table string, records [][]string, skipfirst bool) (sql.Result, error) {
	// read csv file and build sql insert query
	var querystr strings.Builder
	querystr.WriteString(fmt.Sprintf(postgres_insert, table))

	// build sql data info
	count := len(records)
	if skipfirst {
		count--
	}
	datalines := make([]string, count)
	shift := 0
	for ridx, record := range records {
		if skipfirst && ridx == 0 {
			shift = -1
			continue
		}
		fields := make([]string, len(record))
		for fidx, field := range record {
			fields[fidx] = fmt.Sprintf("'%s'", field)
		}
		all_fields := strings.Join(fields, ", ")
		datalines[ridx+shift] = fmt.Sprintf("( %s )", all_fields)
	}

	// add csv records to query string
	all_datalines := strings.Join(datalines, ", ")
	ml.Debug(all_datalines)
	querystr.WriteString(all_datalines)
	querystr.WriteString(";\n")

	// execute query
	full_query := querystr.String()
	ml.Debug(full_query)
	return db.Exec(full_query)
}

func main() {
	// process arguments
	argv := os.Args[1:]
	parser := &docopt.Parser{
		HelpHandler: PrintHelpAndExit,
	}
	opts, _ := parser.ParseArgs(help, argv, version)
	connstr, table, csvfile, skipfirst, purge, debug := ProcessArgs(&opts)

	// log options if requested
	ml.PrintDebug = debug
	ml.Debug(opts)

	// get connection string and open connection
	// TODO: decide which driver to use
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		ml.Debug(err)
	}
	defer db.Close()

	// purge the table if requested
	if purge {
		PurgeTable(db, table)
	}

	records, err := ReadCSV(csvfile)
	if err != nil {
		panic(err)
	}

	res, err := PushCSV(db, table, records, skipfirst)
	if err != nil {
		panic(err)
	}
	rows, _ := res.RowsAffected()
	ml.Debug(fmt.Sprintf("Successfully updated %d records.", rows))
}
