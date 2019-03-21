package persistance

import (
	"github.com/pkg/errors"
	// "log"
	// "database/sql"
	// "github.com/eirannejad/pushcsv/internal/cli"
	"github.com/eirannejad/pushcsv/internal/datafile"
	// _ github.com/go-sql-driver/mysql
)

type MySqlWriter struct {
	DatabaseWriter
}

func (w MySqlWriter) Write(tableData *datafile.TableData) (*Result, error) {
	return nil, errors.New("mysql interface not yet implemented")
}
