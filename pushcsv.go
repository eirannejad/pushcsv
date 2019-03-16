package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/docopt/docopt-go"
	_ "github.com/lib/pq"
)

type MyLog struct {
	PrintDebug bool
}

func (m *MyLog) Debug(args ...interface{}) {
	if m.PrintDebug {
		m.Print(args...)
	}
}

func (m *MyLog) Print(args ...interface{}) {
	log.Print(args...)
}

func main() {
	usage := `Push CSV data to a database.

Usage:
	pushcsv <conn> <table> <csv_file> [--map=<field_maps>] [--debug] [--purge]

Options:
	-h --help     Show this screen.
	--version     Show version.
`
	// process arguments
	argv := os.Args[1:]
	opts, _ := docopt.ParseArgs(usage, argv, "1.0")

	// setup logger
	ml := MyLog{}
	ml.PrintDebug, _ = opts.Bool("--debug")
	// log options if requested
	ml.Debug(opts)

	// get connection string and open connection
	// TODO: decide which driver to use
	connStr, _ := opts.String("<conn>")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		ml.Debug(err)
	}

	// get table, csv file, and other options
	table, _ := opts.String("<table>")

	// start building query string
	// purge the table if requested
	// purge, _ := opts.Bool("--purge")
	// if purge {
	// 	db.Exec(fmt.Sprintf(`TRUNCATE TABLE %s`, table))
	// }

	// open csv file
	csv_file, _ := opts.String("<csv_file>")
	csvfile, err := os.Open(csv_file)
	if err != nil {
		log.Fatal(err)
	}

	// read csv file
	csvreader := csv.NewReader(csvfile)
	for {
		record, err := csvreader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(record)
	}

	// add csv records to query string

	// execute query
	var (
		appname    string
		appversion string
	)

	fmt.Println(table)
	rows, err := db.Query(fmt.Sprintf(`SELECT appname, appversion FROM %s`, table))
	if err != nil {
		ml.Debug(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&appname, &appversion)
		if err != nil {
			ml.Debug(err)
		}
		ml.Print(appname, appversion)
	}
	err = rows.Err()
	if err != nil {
		ml.Debug(err)
	}

	defer db.Close()
}
