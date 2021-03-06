package persistance

import (
	"strings"

	"github.com/pkg/errors"
)

type DBBackendName string

const (
	Postgres DBBackendName = "postgres"
	MongoDB  DBBackendName = "mongodb"
	MySql    DBBackendName = "mysql"
	MSSql    DBBackendName = "sqlserver"
	Sqlite   DBBackendName = "sqlite3"
)

type DatabaseConfig struct {
	NeedsHeaders bool
	ConnString   string
	Backend      DBBackendName
	Username     string
	Password     string
}

func NewDatabaseConfig(connString string) (*DatabaseConfig, error) {
	backend, needsheaders, err := parseUri(connString)
	if err != nil {
		return nil, err
	}

	return &DatabaseConfig{
		NeedsHeaders: needsheaders,
		ConnString:   connString,
		Backend:      backend,
	}, nil
}

func parseUri(connString string) (DBBackendName, bool, error) {
	if strings.HasPrefix(connString, "postgres:") {
		return Postgres, false, nil
	} else if strings.HasPrefix(connString, "mongodb:") {
		return MongoDB, true, nil
	} else if strings.HasPrefix(connString, "mysql:") {
		return MySql, false, nil
	} else if strings.HasPrefix(connString, "sqlserver:") {
		return MSSql, false, nil
	} else if strings.HasPrefix(connString, "sqlite3:") {
		return Sqlite, false, nil
	} else {
		return "", false, errors.New("db is not yet supported")
	}
}
