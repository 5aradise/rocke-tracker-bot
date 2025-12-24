package substorage

import (
	root "bot"
	"bot/config"
	model "bot/internal/models"
	rocketleague "bot/internal/models/rocket-league"
	"bot/internal/storage/queries"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteSubscriptionByTelegramID(t *testing.T) {
	soccer2x2 := model.Subscription{
		Players: rocketleague.P2x2,
		Mode:    rocketleague.Soccer,
	}

	db, upDB, downDB, closeDB := root.DBForTests(t)
	defer closeDB()

	t.Run("Existed subscription", func(t *testing.T) {
		assert := assert.New(t)
		upDB()
		defer downDB()
		stor := New(db)

		_, err := queries.New(db).CreateUser(context.Background(), sql.NullInt64{Int64: 69, Valid: true})
		assert.NoError(err)
		_, err = stor.CreateSubscriptionByTelegramID(context.Background(), 69, soccer2x2)
		assert.NoError(err)

		err = stor.DeleteSubscriptionByTelegramID(context.Background(), 69, soccer2x2)
		assert.NoError(err)
	})

	t.Run("Unexisted user", func(t *testing.T) {
		assert := assert.New(t)
		upDB()
		defer downDB()
		stor := New(db)

		err := stor.DeleteSubscriptionByTelegramID(context.Background(), 69, soccer2x2)
		assert.ErrorIs(err, config.ErrNotFound)
	})

	t.Run("Unexisted subscription", func(t *testing.T) {
		assert := assert.New(t)
		upDB()
		defer downDB()
		stor := New(db)

		_, err := queries.New(db).CreateUser(context.Background(), sql.NullInt64{Int64: 69, Valid: true})
		assert.NoError(err)

		err = stor.DeleteSubscriptionByTelegramID(context.Background(), 69, soccer2x2)
		assert.ErrorIs(err, config.ErrNotFound)
	})
}
