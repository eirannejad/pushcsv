package persistance

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/eirannejad/pushcsv/internal/cli"
	"github.com/eirannejad/pushcsv/internal/datafile"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type GenericSQLWriter struct {
	DatabaseWriter
}

func (w GenericSQLWriter) Purge(tableName string) (*Result, error) {
	// open connection
	db, err := openConnection(&w)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// start transaction
	w.Logger.Debug("opening purge transaction")
	tx, beginErr := db.Begin()
	if beginErr != nil {
		return nil, beginErr
	}
	defer tx.Rollback()

	w.Logger.Debug("purging all existing records")
	var purged int
	if !w.DryRun {
		res, pErr := db.Exec(
			fmt.Sprintf("DELETE FROM %s;", tableName))
		if pErr != nil {
			return nil, pErr
		}
		rows, _ := res.RowsAffected()
		purged = int(rows)
	}

	// commit transaction
	w.Logger.Debug("commiting purge transaction")
	if !w.DryRun {
		txnErr := tx.Commit()
		if txnErr != nil {
			return nil, txnErr
		}
	}

	w.Logger.Debug("preparing report")
	return &Result{
		Message: fmt.Sprintf("successfully purged %d records", purged),
	}, nil
}

func (w GenericSQLWriter) Write(tableData *datafile.TableData) (*Result, error) {
	// open connection
	db, err := openConnection(&w)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// start transaction
	w.Logger.Debug("opening transaction")
	tx, beginErr := db.Begin()
	if beginErr != nil {
		return nil, beginErr
	}
	defer tx.Rollback()

	// purge the table if requested
	if w.PurgeBeforeWrite {
		w.Logger.Debug("purging all existing records")
		if !w.DryRun {
			_, pErr := db.Exec(
				fmt.Sprintf("DELETE FROM %s;", tableData.Name))
			if pErr != nil {
				return nil, pErr
			}
		}
	}

	// generate generic sql insert query
	query, qErr := generateQuery(tableData, w.Logger)
	if qErr != nil {
		return nil, qErr
	}

	// run the insert query
	w.Logger.Debug("executing insert query")
	if !w.DryRun {
		_, eErr := db.Exec(query)
		if eErr != nil {
			return nil, eErr
		}
	}

	// vacuum table if requested
	if w.CompactAfterWrite {
		w.Logger.Debug("vacuuming table")
		if !w.DryRun {
			_, pErr := db.Exec(
				fmt.Sprintf("VACUUM %s;", tableData.Name))
			if pErr != nil {
				return nil, pErr
			}
		}
	}

	// commit transaction
	w.Logger.Debug("commiting transaction")
	if !w.DryRun {
		txnErr := tx.Commit()
		if txnErr != nil {
			return nil, txnErr
		}
	}

	w.Logger.Debug("preparing report")
	return &Result{
		Message: fmt.Sprintf(
			"successfully inserted %d records",
			len(tableData.Records)),
	}, nil
}

func openConnection(w *GenericSQLWriter) (*sql.DB, error) {
	// open connection
	w.Logger.Debug(fmt.Sprintf("opening %s connection", w.Config.Backend))
	cleanConnStr := w.Config.ConnString
	if w.Config.Backend == Sqlite || w.Config.Backend == MySql {
		cleanConnStr = strings.Replace(w.Config.ConnString, string(w.Config.Backend)+":", "", 1)
	}
	return sql.Open(string(w.Config.Backend), cleanConnStr)
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
