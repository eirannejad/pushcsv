package persistance

import (
	"database/sql"
	"fmt"
	"log"
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
	w.Logger.Debug("opening postgresql connection")
	db, err := sql.Open("postgres", w.Config.ConnString)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if len(tableData.Records) > 0 {
		// start transaction
		w.Logger.Debug("opening transaction")
		tx, beginErr := db.Begin()
		if beginErr != nil {
			return nil, beginErr
		}
		defer tx.Rollback()

		// purge the table if requested
		if w.Purge && !w.DryRun {
			w.Logger.Debug("truncating table")
			_, eErr := db.Exec(fmt.Sprintf(`TRUNCATE TABLE %s`, tableData.Name))
			if eErr != nil {
				return nil, eErr
			}
		}

		query, qErr := generateQuery(tableData, w.Logger)
		if qErr != nil {
			return nil, qErr
		}
		if !w.DryRun {
			w.Logger.Debug("executing insert query")
			sqlResult, eErr := db.Exec(query)
			if eErr != nil {
				return nil, eErr
			}

			// commit transaction
			w.Logger.Debug("commiting transaction")
			txnErr := tx.Commit()
			if txnErr != nil {
				log.Fatal(txnErr)
			}

			w.Logger.Debug("preparing report")
			rows, _ := sqlResult.RowsAffected()
			return &Result{
				Message: fmt.Sprintf("successfully updated %d records", int(rows)),
			}, nil
		} else {
			w.Logger.Debug("dry run complete")
			return &Result{
				ResultCode: 2,
				Message:    "processed records but no changed were made to db",
			}, nil
		}
	}
	w.Logger.Debug("nothing to write")
	return &Result{
		ResultCode: 1,
		Message:    "no data to write",
	}, nil
}

func generateQuery(tableData *datafile.TableData, logger *cli.Logger) (string, error) {
	// read csv file and build sql insert query
	var querystr strings.Builder

	if tableData.HasHeaders() {
		logger.Debug("generating insert query with headers")
		querystr.WriteString(fmt.Sprintf("INSERT INTO %s %s values ", tableData.Name, ToSql(&tableData.Headers, false)))
	} else {
		logger.Debug("generating insert query with-out headers")
		querystr.WriteString(fmt.Sprintf("INSERT INTO %s values ", tableData.Name))
	}

	// build sql data info
	logger.Debug("building insert query for data")
	datalines := make([]string, 0)
	for _, record := range tableData.Records {
		datalines = append(datalines, ToSql(&record, true))
	}

	// add csv records to query string
	all_datalines := strings.Join(datalines, ", ")
	logger.Trace(all_datalines)
	querystr.WriteString(all_datalines)
	querystr.WriteString(";\n")
	logger.Debug("building query completed")

	// execute query
	full_query := querystr.String()
	logger.Trace(full_query)
	return full_query, nil
}
