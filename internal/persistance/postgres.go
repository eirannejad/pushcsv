package persistance

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/eirannejad/pushcsv/internal/cli"
	"github.com/eirannejad/pushcsv/internal/datafile"
	_ "github.com/lib/pq"
)

type PostgresWriter struct {
	DatabaseWriter
}

func (w PostgresWriter) Write(tableData *datafile.TableData) (*Result, error) {
	// open connection
	db, err := sql.Open("postgres", w.ConnectionUri)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// purge the table if requested
	if w.Purge && !w.DryRun {
		_, eErr := db.Exec(fmt.Sprintf(`TRUNCATE TABLE %s`, tableData.Name))
		if eErr != nil {
			return nil, eErr
		}
	}

	if len(tableData.Records) > 0 {
		query, qErr := generateQuery(tableData, w.Logger)
		if qErr != nil {
			return nil, qErr
		}
		if !w.DryRun {
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
	}
	return &Result{
		Message: "No data to write",
		Count:   0,
	}, nil
}

func generateQuery(tableData *datafile.TableData, logger *cli.Logger) (string, error) {
	// read csv file and build sql insert query
	var querystr strings.Builder

	if tableData.HasHeaders() {
		querystr.WriteString(fmt.Sprintf("INSERT INTO %s %s values ", tableData.Name, ToSql(&tableData.Headers)))
	} else {
		querystr.WriteString(fmt.Sprintf("INSERT INTO %s values ", tableData.Name))
	}

	// build sql data info
	datalines := make([]string, 0)
	for _, record := range tableData.Records {
		datalines = append(datalines, ToSql(&record))
	}

	// add csv records to query string
	all_datalines := strings.Join(datalines, ", ")
	logger.Trace(all_datalines)
	querystr.WriteString(all_datalines)
	querystr.WriteString(";\n")

	// execute query
	full_query := querystr.String()
	logger.Trace(full_query)
	return full_query, nil
}
