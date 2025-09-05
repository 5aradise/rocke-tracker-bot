package subStorage

import (
	"bot/config"
	model "bot/internal/models"
	"bot/internal/storage/adapter"
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

func (s Storage) CreateSubscriptionByTelegramID(ctx context.Context, tgID int64, sub model.Subscription) (int64, error) {
	subID, err := queries.New(s.db).CreateSubscriptionByTelegramID(ctx, queries.CreateSubscriptionByTelegramIDParams{
		TelegramID: sql.NullInt64{
			Int64: tgID,
			Valid: true,
		},
		Players: adapter.PlayersToDB(sub.Players),
		Mode:    adapter.ModeToDB(sub.Mode),
	})
	if err != nil {
		switch sqlite.ErrorType(err) {
		case sqlite.CONSTRAINT_UNIQUE:
			return 0, config.ErrUniqueConstraint
		case sqlite.CONSTRAINT_NOTNULL:
			return 0, config.ErrNotFound
		default:
			return 0, err
		}
	}

	return subID, nil
}

func (s Storage) ListTelegramIDsBySubscription(ctx context.Context, sub model.Subscription) ([]int64, error) {
	tgIDs, err := queries.New(s.db).ListTelegramIDsBySubscription(ctx, queries.ListTelegramIDsBySubscriptionParams{
		Players: adapter.PlayersToDB(sub.Players),
		Mode:    adapter.ModeToDB(sub.Mode),
	})
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0, len(tgIDs))
	for _, id := range tgIDs {
		ids = append(ids, id.Int64)
	}
	return ids, nil
}
