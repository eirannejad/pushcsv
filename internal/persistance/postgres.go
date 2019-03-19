package persistance

import (
	"database/sql"
	"fmt"
	"github.com/eirannejad/pushcsv/internal/cli"
	"github.com/eirannejad/pushcsv/internal/csv"
	_ "github.com/lib/pq"
	"strings"
)

type PostgresWriter struct {
	DatabaseWriter
}

func (w PostgresWriter) Write(csvData *csv.CsvData) (*Result, error) {
	// open connection
	db, err := sql.Open("postgres", w.Options.ConnString)
	if err != nil {
		w.Logger.Debug(err)
	}
	defer db.Close()

	// purge the table if requested
	if w.Options.Purge {
		db.Exec(fmt.Sprintf(`TRUNCATE TABLE %s`, w.Options.Table))
	}

	if len(csvData.Records) > 0 {
		query, qErr := GenerateQuery(w.Logger, w.Options, csvData)
		if qErr != nil {
			return nil, qErr
		}
		sqlResult, eErr := db.Exec(query)
		if eErr != nil {
			return nil, eErr
		}
		rows, _ := sqlResult.RowsAffected()
		return &Result{
			Message: "Data Writen",
			Count:   int(rows),
		}, nil
	}
	return &Result{
		Message: "No data to write",
		Count:   0,
	}, nil
}

func GenerateQuery(logger *cli.Logger, options *cli.Options, csvData *csv.CsvData) (string, error) {
	// read csv file and build sql insert query
	var querystr strings.Builder

	if len(options.AttrMap) > 0 {
		columns := fmt.Sprintf("( %s )", strings.Join(options.AttrMap, ","))
		querystr.WriteString(fmt.Sprintf("INSERT INTO %s %s values ", options.Table, columns))
	} else {
		querystr.WriteString(fmt.Sprintf("INSERT INTO %s values ", options.Table))
	}

	// build sql data info
	count := len(csvData.Records)
	datalines := make([]string, count)
	for ridx, record := range csvData.Records {
		fields := make([]string, len(record))
		for fidx, field := range record {
			fields[fidx] = fmt.Sprintf("'%s'", field)
		}
		all_fields := strings.Join(fields, ", ")
		datalines[ridx] = fmt.Sprintf("( %s )", all_fields)
	}

	// add csv records to query string
	all_datalines := strings.Join(datalines, ", ")
	logger.Debug(all_datalines)
	querystr.WriteString(all_datalines)
	querystr.WriteString(";\n")

	// execute query
	full_query := querystr.String()
	logger.Debug(full_query)
	return full_query, nil
}
