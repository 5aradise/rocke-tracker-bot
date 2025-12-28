package sqlite

import (
	"errors"

	"modernc.org/sqlite"
)

type ErrType int

const (
	CONSTRAINT_NOTNULL = 1299
	CONSTRAINT_UNIQUE  = 2067
)

// func IsUniqueConstraintError(err error) bool {
// 	var sqliteError *sqlite.Error
// 	if errors.As(err, &sqliteError) {
// 		return sqliteError.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE
// 	}
// 	panic("not modernc.org/sqlite/lib error")
// }

func ErrorType(err error) ErrType {
	var sqliteError *sqlite.Error
	if errors.As(err, &sqliteError) {
		return ErrType(sqliteError.Code())
	}
	panic("not modernc.org/sqlite/lib error")
}
