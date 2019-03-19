package persistance

import (
	"github.com/pkg/errors"
	"strings"
)

type DBBackendName string

const (
	Postgres DBBackendName = "postgres"
	Sqlite   DBBackendName = "sqlite"
	MongoDB  DBBackendName = "mongodb"
)

type DatabaseConfig struct {
	ConnString string
	Backend    DBBackendName
}

func NewDatabaseConfig(connString string) (*DatabaseConfig, error) {
	backend, err := parseUri(connString)
	if err != nil {
		return nil, err
	}
	return &DatabaseConfig{
		ConnString: connString,
		Backend:    backend,
	}, nil
}

func parseUri(connString string) (DBBackendName, error) {
	var db DBBackendName
	if strings.HasPrefix(connString, "postgres:") {
		db = Postgres
	} else if strings.HasPrefix(connString, "sqlite:") {
		db = Sqlite
	} else if strings.HasPrefix(connString, "mongodb:") {
		db = MongoDB
	} else {
		return "", errors.New("unsupported db")
	}
	return db, nil
}
