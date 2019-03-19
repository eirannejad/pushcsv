package persistance

import (
	"github.com/pkg/errors"
	"strings"
)

type DBBackend string

const (
	Postgres DBBackend = "postgres"
	Sqlite   DBBackend = "sqlite"
	MongoDB  DBBackend = "mongodb"
)

type DatabaseConfig struct {
	ConnString string
	Backend    DBBackend
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

func parseUri(connString string) (DBBackend, error) {
	var db DBBackend
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
