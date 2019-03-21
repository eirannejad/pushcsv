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
	Sqlite   DBBackendName = "sqlite"
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

	// TODO: implement auth from env var
	// var username, password string
	// authinfo := os.Getenv("PUSHCSVAUTH")
	// if authinfo != "" {
	// 	parts := strings.Split(authinfo, ":")
	// 	if len(parts) == 2 {
	// 		username = parts[0]
	// 		password = parts[1]
	// 	}
	// }

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
		return MySql, true, nil
	} else if strings.HasPrefix(connString, "sqlite:") {
		return Sqlite, false, nil
	} else {
		return "", false, errors.New("db is not yet supported")
	}
}
