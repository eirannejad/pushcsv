package persistance

import (
	// "errors"
	"github.com/eirannejad/pushcsv/internal/cli"
	"github.com/eirannejad/pushcsv/internal/datafile"
	_ "github.com/go-sql-driver/mysql"
	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	// "strings"
)

// ErroCodes
// 0: All OK
// 1: No data to write
// 2: data is available but did not get pushed under dry run
// 3: headers are required
type Result struct {
	ResultCode int
	Message    string
}

type DatabaseWriter struct {
	Config *DatabaseConfig
	Purge  bool
	DryRun bool
	Logger *cli.Logger
}

type Writer interface {
	// Ensure whatever is return has Write Method
	Write(*datafile.TableData) (*Result, error)
}

func NewWriter(dbConfig *DatabaseConfig, options *cli.Options, logger *cli.Logger) (Writer, error) {
	w := &DatabaseWriter{
		Config: dbConfig,
		Purge:  options.Purge,
		DryRun: options.DryRun,
		Logger: logger,
	}
	if dbConfig.Backend == Postgres {
		return PostgresWriter{*w}, nil
	} else if dbConfig.Backend == MongoDB {
		return MongoDBWriter{*w}, nil
	} else if dbConfig.Backend == MySql {
		return MySqlWriter{*w}, nil
	} else if dbConfig.Backend == Sqlite {
		return SqliteWriter{*w}, nil
	}
	// ... other writers

	panic("should not get here")
}
