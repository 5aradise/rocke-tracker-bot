package sqlite

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

const (
	driverName             = "sqlite"
	fkConstraintActivation = "?_pragma=foreign_keys(1)"
)

func New(file string) (*sql.DB, error) {
	conn, err := sql.Open(driverName, file+fkConstraintActivation)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	return conn, nil
}
