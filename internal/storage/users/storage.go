package userStorage

import (
	"bot/config"
	"bot/internal/storage/queries"
	"bot/pkg/sqlite"
	"context"
	"database/sql"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) Storage {
	return Storage{
		db: db,
	}
}

func (s Storage) CreateUser(ctx context.Context, tgID int64) (int64, error) {
	userID, err := queries.New(s.db).CreateUser(ctx, sql.NullInt64{
		Int64: tgID,
		Valid: true,
	})
	if err != nil {
		switch sqlite.ErrorType(err) {
		case sqlite.CONSTRAINT_UNIQUE:
			return 0, config.ErrUniqueConstraint
		default:
			return 0, err
		}
	}

	return userID, nil
}
