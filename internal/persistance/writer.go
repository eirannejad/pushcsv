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

type Result struct {
	Message string
	Count   int
	// Records: ?
}

type DatabaseWriter struct {
	ConnectionUri string
	Purge         bool
	DryRun        bool
	Logger        *cli.Logger
}

type Writer interface {
	// Ensure whatever is return has Write Method
	Write(*datafile.TableData) (*Result, error)
}

func NewWriter(dbConfig *DatabaseConfig, options *cli.Options, logger *cli.Logger) (Writer, error) {
	w := &DatabaseWriter{
		ConnectionUri: options.ConnString,
		Purge:         options.Purge,
		DryRun:        options.DryRun,
		Logger:        logger,
	}
	if dbConfig.Backend == Postgres {
		return PostgresWriter{*w}, nil
	} else if dbConfig.Backend == MongoDB {
		return MongoDBWriter{*w}, nil
	} else if dbConfig.Backend == Sqlite {
		return SqliteWriter{*w}, nil
	}
	// ... other writers

	panic("Should never get here")
}
