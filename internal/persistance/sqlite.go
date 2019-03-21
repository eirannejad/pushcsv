package persistance

import (
	"github.com/pkg/errors"
	// "log"
	// "database/sql"
	// "github.com/eirannejad/pushcsv/internal/cli"
	"github.com/eirannejad/pushcsv/internal/datafile"
	// _ "github.com/mattn/go-sqlite3"
)

type SqliteWriter struct {
	DatabaseWriter
}

func (w SqliteWriter) Write(tableData *datafile.TableData) (*Result, error) {
	return nil, errors.New("sqlite interface not yet implemented.")
}
