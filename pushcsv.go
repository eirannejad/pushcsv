package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/docopt/docopt-go"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	// _ "github.com/mattn/go-sqlite3"
)

type MyLog struct {
	PrintDebug bool
}

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

func ParseArgs(opts *docopt.Opts) (string, string, string, bool, bool, []string, bool) {
	conn_string, _ := opts.String("<db_uri>")
	table, _ := opts.String("<table>")
	csvfile, _ := opts.String("<csv_file>")
	has_headers, _ := opts.Bool("--headers")
	purge, _ := opts.Bool("--purge")
	attrmap, _ := opts.String("--map")
	debug, _ := opts.Bool("--debug")
	if attrmap != "" {
		return conn_string, table, csvfile, has_headers, purge, strings.Split(attrmap, ";"), debug
	} else {
		return conn_string, table, csvfile, has_headers, purge, make([]string, 0), debug
	}
}

func ReadCSV(csvfile string, has_headers bool) ([][]string, []string) {
	// open csv file
	csvfilehndlr, err := os.Open(csvfile)
	if err != nil {
		panic(err)
	}

	csvreader := csv.NewReader(csvfilehndlr)
	// csvreader.LazyQuotes = true
	records, err := csvreader.ReadAll()
	if err != nil {
		panic(err)
	}
	if has_headers {
		return records[1:], records[0]
	} else {
		return records, nil
	}
}

func PushSQLRecords(db *sql.DB, table string, records [][]string, headers []string, attrs []string) (sql.Result, error) {
	// read csv file and build sql insert query
	var querystr strings.Builder

	if len(attrs) > 0 {
		columns := fmt.Sprintf("( %s )", strings.Join(attrs, ","))
		querystr.WriteString(fmt.Sprintf("INSERT INTO %s %s values ", table, columns))
	} else {
		querystr.WriteString(fmt.Sprintf("INSERT INTO %s values ", table))
	}

	// build sql data info
	count := len(records)
	datalines := make([]string, count)
	for ridx, record := range records {
		fields := make([]string, len(record))
		for fidx, field := range record {
			fields[fidx] = fmt.Sprintf("'%s'", field)
		}
		all_fields := strings.Join(fields, ", ")
		datalines[ridx] = fmt.Sprintf("( %s )", all_fields)
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

func PushMongoDB(connstr string, table string, docs [][]string, purge bool, headers []string, attrs []string) (int, error) {
	// parse and grab database name from uri
	dialinfo, err := mgo.ParseURL(connstr)
	if err != nil {
		panic(err)
	}
	dbname := dialinfo.Database

	// connect to db engine
	session, err := mgo.Dial(connstr)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(dbname).C(table)

	// purge the collection if requested
	if purge {
		c.RemoveAll(bson.M{})
	}

	if len(docs) > 0 {
		if len(attrs) == 0 {
			if len(headers) == 0 {
				log.Fatal("`--map` must be specified when pushing a csv with no headers.")
				return 0, nil
			} else {
				attrs = headers
			}
		}

		// build sql data info
		bulkop := c.Bulk()
		for _, record := range docs {
			map_obj := make(map[string]string)
			for fidx, field := range record {
				map_obj[attrs[fidx]] = field
			}
			bulkop.Insert(map_obj)
		}
		res, err := bulkop.Run()
		if err != nil {
			log.Fatal(err)
		}
		return res.Modified, nil
	}
	return 0, nil
}

func PushSqlite(connstr string, table string, records [][]string, purge bool, headers []string, attrs []string) (int, error) {
	log.Fatal("sqlite interface not yet implemented.")
	return 0, nil
}

func PushPostgres(connstr string, table string, records [][]string, purge bool, headers []string, attrs []string) (int, error) {
	// open connection
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		ml.Debug(err)
	}
	defer db.Close()

	// purge the table if requested
	if purge {
		db.Exec(fmt.Sprintf(`TRUNCATE TABLE %s`, table))
	}

	if len(records) > 0 {
		res, err := PushSQLRecords(db, table, records, headers, attrs)
		if err != nil {
			panic(err)
		}
		rows, _ := res.RowsAffected()
		return int(rows), nil
	}
	return 0, nil
}

func main() {
	// process arguments
	argv := os.Args[1:]
	parser := &docopt.Parser{
		HelpHandler: PrintHelpAndExit,
	}
	opts, _ := parser.ParseArgs(help, argv, version)
	connstr, table, csvfile, has_headers, purge, attrs, debug := ParseArgs(&opts)

	// log options if requested
	ml.PrintDebug = debug
	ml.Debug(opts)

	// read csv (handles panics)
	var records [][]string = [][]string{}
	var headers []string = []string{}
	if csvfile != "" {
		records, headers = ReadCSV(csvfile, has_headers)
	}

	// check connection string and determine target db driver
	postgres := strings.HasPrefix(connstr, "postgres:")
	sqlite := strings.HasPrefix(connstr, "sqlite:")
	mongodb := strings.HasPrefix(connstr, "mongodb:")

	var pusherr error
	var modified int = 0
	if postgres {
		modified, pusherr = PushPostgres(connstr, table, records, purge, headers, attrs)
	} else if sqlite {
		modified, pusherr = PushSqlite(connstr, table, records, purge, headers, attrs)
	} else if mongodb {
		modified, pusherr = PushMongoDB(connstr, table, records, purge, headers, attrs)
	}

	if pusherr != nil {
		panic(pusherr)
	}
	ml.Debug(fmt.Sprintf("Successfully updated %d records.", modified))
}
