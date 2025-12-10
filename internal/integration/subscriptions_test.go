package integration

import (
	root "bot"
	"bot/config"
	model "bot/internal/models"
	rocketleague "bot/internal/models/rocket-league"
	subservice "bot/internal/services/subscriptions"
	userservice "bot/internal/services/users"
	substorage "bot/internal/storage/subscriptions"
	userstorage "bot/internal/storage/users"
	"bot/pkg/sqlite"
	"context"
	"testing"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscriptionsServiceAndStorage(t *testing.T) {
	require := require.New(t)
	soccer2x2 := model.Subscription{
		Players: rocketleague.P2x2,
		Mode:    rocketleague.Soccer,
	}
	soccer3x3 := model.Subscription{
		Players: rocketleague.P3x3,
		Mode:    rocketleague.Soccer,
	}
	pentathlon3x3 := model.Subscription{
		Players: rocketleague.P3x3,
		Mode:    rocketleague.Pentathlon,
	}
	db, err := sqlite.New(":memory:")
	if err != nil {
		require.NoError(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			require.NoError(err)
		}
	}()
	goose.SetBaseFS(root.ForIntegrationMigrations)
	if err := goose.SetDialect("sqlite"); err != nil {
		require.NoError(err)
	}
	if err := goose.Up(db, "sql/schema"); err != nil {
		require.NoError(err)
	}

	t.Run("Simple", func(t *testing.T) {
		// Arrange
		assert := assert.New(t)
		// db
		if err := goose.Up(db, "sql/schema"); err != nil {
			require.NoError(err)
		}
		defer func() {
			if err := goose.Down(db, "sql/schema"); err != nil {
				require.NoError(err)
			}
		}()
		// storages
		userStor := userstorage.New(db)
		subStor := substorage.New(db)
		// services
		userServ := userservice.New(userStor)
		subServ := subservice.New(nil, subStor)

		// Act
		_, serr := userServ.CreateUser(context.Background(), 1)
		assert.Equal(config.NilError, serr)
		_, serr = subServ.SubscribeByTelegram(context.Background(), 1, soccer2x2)
		assert.Equal(config.NilError, serr)
		_, serr = subServ.SubscribeByTelegram(context.Background(), 1, pentathlon3x3)
		assert.Equal(config.NilError, serr)
		subs, serr := subServ.ListTelegramUserSubscriptions(context.Background(), 1)
		assert.Equal(config.NilError, serr)

		// Assert
		assert.Contains(subs, soccer2x2)
		assert.Contains(subs, pentathlon3x3)
		assert.NotContains(subs, soccer3x3)
	})

	t.Run("Multiple users", func(t *testing.T) {
		// Arrange
		assert := assert.New(t)
		// db
		if err := goose.Up(db, "sql/schema"); err != nil {
			require.NoError(err)
		}
		defer func() {
			if err := goose.Down(db, "sql/schema"); err != nil {
				require.NoError(err)
			}
		}()
		// storages
		userStor := userstorage.New(db)
		subStor := substorage.New(db)
		// services
		userServ := userservice.New(userStor)
		subServ := subservice.New(nil, subStor)

		// Act
		_, serr := userServ.CreateUser(context.Background(), 1)
		assert.Equal(config.NilError, serr)
		_, serr = userServ.CreateUser(context.Background(), 69)
		assert.Equal(config.NilError, serr)
		_, serr = userServ.CreateUser(context.Background(), 420)
		assert.Equal(config.NilError, serr)
		_, serr = subServ.SubscribeByTelegram(context.Background(), 1, soccer2x2)
		assert.Equal(config.NilError, serr)
		_, serr = subServ.SubscribeByTelegram(context.Background(), 1, pentathlon3x3)
		assert.Equal(config.NilError, serr)
		_, serr = subServ.SubscribeByTelegram(context.Background(), 69, soccer2x2)
		assert.Equal(config.NilError, serr)
		_, serr = subServ.SubscribeByTelegram(context.Background(), 69, pentathlon3x3)
		assert.Equal(config.NilError, serr)
		_, serr = subServ.SubscribeByTelegram(context.Background(), 69, soccer3x3)
		assert.Equal(config.NilError, serr)
		_, serr = subServ.SubscribeByTelegram(context.Background(), 420, soccer2x2)
		assert.Equal(config.NilError, serr)

		// Assert
		subs, serr := subServ.ListTelegramUserSubscriptions(context.Background(), 1)
		assert.Equal(config.NilError, serr)
		assert.Contains(subs, soccer2x2)
		assert.Contains(subs, pentathlon3x3)
		assert.NotContains(subs, soccer3x3)

		subs, serr = subServ.ListTelegramUserSubscriptions(context.Background(), 69)
		assert.Equal(config.NilError, serr)
		assert.Contains(subs, soccer2x2)
		assert.Contains(subs, pentathlon3x3)
		assert.Contains(subs, soccer3x3)

		subs, serr = subServ.ListTelegramUserSubscriptions(context.Background(), 420)
		assert.Equal(config.NilError, serr)
		assert.Contains(subs, soccer2x2)
		assert.NotContains(subs, pentathlon3x3)
		assert.NotContains(subs, soccer3x3)
	})

	t.Run("Unexisted user subscription", func(t *testing.T) {
		// Arrange
		assert := assert.New(t)
		// db
		if err := goose.Up(db, "sql/schema"); err != nil {
			require.NoError(err)
		}
		defer func() {
			if err := goose.Down(db, "sql/schema"); err != nil {
				require.NoError(err)
			}
		}()
		// storages
		userStor := userstorage.New(db)
		subStor := substorage.New(db)
		// services
		userServ := userservice.New(userStor)
		subServ := subservice.New(nil, subStor)

		// Act
		_, serr := userServ.CreateUser(context.Background(), 1)
		assert.Equal(config.NilError, serr)
		_, serr = subServ.SubscribeByTelegram(context.Background(), 3, soccer2x2)

		// Assert
		assert.Equal(config.CodeUserWithTgIDNotExist, serr.Code)
	})

	t.Run("Same subscription", func(t *testing.T) {
		// Arrange
		assert := assert.New(t)
		// db
		if err := goose.Up(db, "sql/schema"); err != nil {
			require.NoError(err)
		}
		defer func() {
			if err := goose.Down(db, "sql/schema"); err != nil {
				require.NoError(err)
			}
		}()
		// storages
		userStor := userstorage.New(db)
		subStor := substorage.New(db)
		// services
		userServ := userservice.New(userStor)
		subServ := subservice.New(nil, subStor)

		// Act
		_, serr := userServ.CreateUser(context.Background(), 1)
		assert.Equal(config.NilError, serr)
		_, serr = subServ.SubscribeByTelegram(context.Background(), 1, soccer2x2)
		assert.Equal(config.NilError, serr)
		_, serr = subServ.SubscribeByTelegram(context.Background(), 1, soccer2x2)

		// Assert
		assert.Equal(config.CodeUserHasSub, serr.Code)
	})

	t.Run("Subscriptions of unexisted user should just return empty slice", func(t *testing.T) {
		// Arrange
		assert := assert.New(t)
		// db
		if err := goose.Up(db, "sql/schema"); err != nil {
			require.NoError(err)
		}
		defer func() {
			if err := goose.Down(db, "sql/schema"); err != nil {
				require.NoError(err)
			}
		}()
		// storages
		userStor := userstorage.New(db)
		subStor := substorage.New(db)
		// services
		userServ := userservice.New(userStor)
		subServ := subservice.New(nil, subStor)

		// Act
		_, serr := userServ.CreateUser(context.Background(), 1)
		assert.Equal(config.NilError, serr)
		_, serr = subServ.SubscribeByTelegram(context.Background(), 1, soccer2x2)
		assert.Equal(config.NilError, serr)
		_, serr = subServ.SubscribeByTelegram(context.Background(), 1, pentathlon3x3)
		assert.Equal(config.NilError, serr)
		subs, serr := subServ.ListTelegramUserSubscriptions(context.Background(), 3)

		// Assert
		assert.Empty(subs)
		assert.Equal(config.NilError, serr)
	})
}
