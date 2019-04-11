package persistance

import (
	"github.com/eirannejad/pushcsv/internal/cli"
	"github.com/eirannejad/pushcsv/internal/datafile"
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
	Config           *DatabaseConfig
	PurgeBeforeWrite bool
	DryRun           bool
	Logger           *cli.Logger
}

type Writer interface {
	// Ensure whatever is return has Write Method
	Write(*datafile.TableData) (*Result, error)
	Purge(string) (*Result, error)
}

func NewWriter(dbConfig *DatabaseConfig, options *cli.Options, logger *cli.Logger) (Writer, error) {
	w := &DatabaseWriter{
		Config:           dbConfig,
		PurgeBeforeWrite: options.Purge,
		DryRun:           options.DryRun,
		Logger:           logger,
	}
	if dbConfig.Backend == Postgres {
		return GenericSQLWriter{*w}, nil
	} else if dbConfig.Backend == MongoDB {
		return MongoDBWriter{*w}, nil
	} else if dbConfig.Backend == MySql {
		return GenericSQLWriter{*w}, nil
	} else if dbConfig.Backend == MSSql {
		return GenericSQLWriter{*w}, nil
	} else if dbConfig.Backend == Sqlite {
		return GenericSQLWriter{*w}, nil
	}
	// ... other writers

	panic("should not get here")
}
