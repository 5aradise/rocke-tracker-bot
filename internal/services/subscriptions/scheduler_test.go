package subservice

import (
	"bot/config"
	model "bot/internal/models"
	rocketleague "bot/internal/models/rocket-league"
	"context"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSubscribeByTelegram(t *testing.T) {
	soccer2x2 := model.Subscription{
		Players: rocketleague.P2x2,
		Mode:    rocketleague.Soccer,
	}
	pentathlon3x3 := model.Subscription{
		Players: rocketleague.P3x3,
		Mode:    rocketleague.Pentathlon,
	}

	t.Run("Unexisted user subscription", func(t *testing.T) {
		assert := assert.New(t)
		s := New(&mockRocketLeagueAPI{}, &mockSubStorage{
			returnSubID:  69,
			unexistedIDs: []int64{1337},
		})

		id, err := s.SubscribeByTelegram(context.Background(), 1337, soccer2x2)

		assert.Zero(id)
		assert.Equal(config.CodeUserWithTgIDNotExist, err.Code)
	})

	t.Run("User already has this subscription", func(t *testing.T) {
		assert := assert.New(t)
		s := New(&mockRocketLeagueAPI{}, &mockSubStorage{
			returnSubID: 69,
		})

		id, err := s.SubscribeByTelegram(context.Background(), 1337, soccer2x2)
		assert.EqualValues(69, id)
		assert.Equal(config.NilError, err)
		id, err = s.SubscribeByTelegram(context.Background(), 1337, pentathlon3x3)
		assert.EqualValues(69, id)
		assert.Equal(config.NilError, err)
		id, err = s.SubscribeByTelegram(context.Background(), 1337, soccer2x2)

		assert.Zero(id)
		assert.Equal(config.CodeUserHasSub, err.Code)
	})
}

func TestRunNotifications(t *testing.T) {
	updateDelay = 24 * time.Hour
	tournamentDelay = 1 * time.Second
	timeForDB = 100 * time.Nanosecond

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

	t.Run("In time notifications", func(t *testing.T) {
		soccer2x2 := model.Subscription{
			Players: rocketleague.P2x2,
			Mode:    rocketleague.Soccer,
		}
		timeForRegistration := time.Second
		tournamentStart := time.Now().Add(timeForRegistration + tournamentDelay)

		assert := assert.New(t)
		s := New(
			&mockRocketLeagueAPI{
				tours: []model.Tournament{
					{
						Type:   soccer2x2,
						Starts: tournamentStart,
					},
				},
			},
			&mockSubStorage{
				returnSubID: 69,
			},
		)
		ntfCh := make(chan model.TgNotification)

		go s.RunNotifications(ntfCh)
		subID, err := s.SubscribeByTelegram(context.Background(), 1, soccer2x2)
		assert.EqualValues(69, subID)
		assert.Equal(config.NilError, err)
		subID, err = s.SubscribeByTelegram(context.Background(), 69, soccer2x2)
		assert.EqualValues(69, subID)
		assert.Equal(config.NilError, err)
		subID, err = s.SubscribeByTelegram(context.Background(), 420, soccer2x2)
		assert.EqualValues(69, subID)
		assert.Equal(config.NilError, err)
		ntf := <-ntfCh

		assertNowIs(t, tournamentStart.Add(-tournamentDelay))
		assert.Equal(soccer2x2, ntf.Tournament)
		assert.Contains(ntf.IDs, int64(1))
		assert.Contains(ntf.IDs, int64(69))
		assert.Contains(ntf.IDs, int64(420))
	})

	t.Run("Notification only subscribed", func(t *testing.T) {
		timeForRegistration := time.Second
		tournamentStart := time.Now().Add(timeForRegistration + tournamentDelay)

		assert := assert.New(t)
		s := New(
			&mockRocketLeagueAPI{
				tours: []model.Tournament{
					{
						Type:   soccer2x2,
						Starts: tournamentStart,
					},
				},
			},
			&mockSubStorage{
				returnSubID: 69,
			},
		)
		ntfCh := make(chan model.TgNotification)

		go s.RunNotifications(ntfCh)
		subID, err := s.SubscribeByTelegram(context.Background(), 1, soccer2x2)
		assert.EqualValues(69, subID)
		assert.Equal(config.NilError, err)
		subID, err = s.SubscribeByTelegram(context.Background(), 69, pentathlon3x3)
		assert.EqualValues(69, subID)
		assert.Equal(config.NilError, err)
		subID, err = s.SubscribeByTelegram(context.Background(), 420, pentathlon3x3)
		assert.EqualValues(69, subID)
		assert.Equal(config.NilError, err)
		subID, err = s.SubscribeByTelegram(context.Background(), 1337, soccer2x2)
		assert.EqualValues(69, subID)
		assert.Equal(config.NilError, err)
		ntf := <-ntfCh

		assertNowIs(t, tournamentStart.Add(-tournamentDelay))
		assert.Equal(soccer2x2, ntf.Tournament)
		assert.Contains(ntf.IDs, int64(1))
		assert.NotContains(ntf.IDs, int64(69))
		assert.NotContains(ntf.IDs, int64(420))
		assert.Contains(ntf.IDs, int64(1337))
	})

	t.Run("Notification for several tournaments", func(t *testing.T) {
		timeForRegistration := time.Second
		var (
			soccer2x2Start     = time.Now().Add(0*time.Second + timeForRegistration + tournamentDelay)
			soccer3x3Start     = time.Now().Add(1*time.Second + timeForRegistration + tournamentDelay)
			pentathlon3x3Start = time.Now().Add(2*time.Second + timeForRegistration + tournamentDelay)
		)

		assert := assert.New(t)
		s := New(
			&mockRocketLeagueAPI{
				tours: []model.Tournament{
					{
						Type:   soccer2x2,
						Starts: soccer2x2Start,
					},
					{
						Type:   soccer3x3,
						Starts: soccer3x3Start,
					},
					{
						Type:   pentathlon3x3,
						Starts: pentathlon3x3Start,
					},
				},
			},
			&mockSubStorage{
				returnSubID: 69,
			},
		)
		ntfCh := make(chan model.TgNotification)

		go s.RunNotifications(ntfCh)

		subID, err := s.SubscribeByTelegram(context.Background(), 1, soccer2x2)
		assert.EqualValues(69, subID)
		assert.Equal(config.NilError, err)
		subID, err = s.SubscribeByTelegram(context.Background(), 69, soccer2x2)
		assert.EqualValues(69, subID)
		assert.Equal(config.NilError, err)
		subID, err = s.SubscribeByTelegram(context.Background(), 420, soccer2x2)
		assert.EqualValues(69, subID)
		assert.Equal(config.NilError, err)

		subID, err = s.SubscribeByTelegram(context.Background(), 1, soccer3x3)
		assert.EqualValues(69, subID)
		assert.Equal(config.NilError, err)
		subID, err = s.SubscribeByTelegram(context.Background(), 420, soccer3x3)
		assert.EqualValues(69, subID)
		assert.Equal(config.NilError, err)
		subID, err = s.SubscribeByTelegram(context.Background(), 1337, soccer3x3)
		assert.EqualValues(69, subID)
		assert.Equal(config.NilError, err)

		subID, err = s.SubscribeByTelegram(context.Background(), 69, pentathlon3x3)
		assert.EqualValues(69, subID)
		assert.Equal(config.NilError, err)
		subID, err = s.SubscribeByTelegram(context.Background(), 1337, pentathlon3x3)
		assert.EqualValues(69, subID)
		assert.Equal(config.NilError, err)

		ntf := <-ntfCh
		assertNowIs(t, soccer2x2Start.Add(-tournamentDelay))
		assert.Equal(soccer2x2, ntf.Tournament)
		assert.Contains(ntf.IDs, int64(1))
		assert.Contains(ntf.IDs, int64(69))
		assert.Contains(ntf.IDs, int64(420))
		assert.NotContains(ntf.IDs, int64(1337))

		ntf = <-ntfCh
		assertNowIs(t, soccer3x3Start.Add(-tournamentDelay))
		assert.Equal(soccer3x3, ntf.Tournament)
		assert.Contains(ntf.IDs, int64(1))
		assert.NotContains(ntf.IDs, int64(69))
		assert.Contains(ntf.IDs, int64(420))
		assert.Contains(ntf.IDs, int64(1337))

		ntf = <-ntfCh
		assertNowIs(t, pentathlon3x3Start.Add(-tournamentDelay))
		assert.Equal(pentathlon3x3, ntf.Tournament)
		assert.NotContains(ntf.IDs, int64(1))
		assert.Contains(ntf.IDs, int64(69))
		assert.NotContains(ntf.IDs, int64(420))
		assert.Contains(ntf.IDs, int64(1337))
	})
}

func TestAtFunc(t *testing.T) {
	after100ms := time.Now().Add(100 * time.Millisecond)
	afterOneSec := time.Now().Add(time.Second)
	afterSomeTime := time.Now().Add(3333 * time.Millisecond)

	atFunc(after100ms, func() {
		assertNowIs(t, after100ms)
	})

	atFunc(afterOneSec, func() {
		assertNowIs(t, afterOneSec)
	})

	atFunc(afterSomeTime, func() {
		assertNowIs(t, afterSomeTime)
	})
}

type mockRocketLeagueAPI struct {
	mu    sync.Mutex
	tours []model.Tournament
}

func (m *mockRocketLeagueAPI) Tournaments() ([]model.Tournament, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return slices.Clone(m.tours), nil
}

type mockSubStorage struct {
	returnSubID  int64
	mu           sync.Mutex
	idsBySub     map[model.Subscription][]int64
	unexistedIDs []int64
}

func (m *mockSubStorage) CreateSubscriptionByTelegramID(ctx context.Context, tgID int64, sub model.Subscription) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if slices.Contains(m.unexistedIDs, tgID) {
		return 0, config.ErrNotFound
	}

	if m.idsBySub == nil {
		m.idsBySub = make(map[model.Subscription][]int64)
	}

	subIDs := m.idsBySub[sub]
	if slices.Contains(subIDs, tgID) {
		return 0, config.ErrUniqueConstraint
	}

	subIDs = append(subIDs, tgID)
	m.idsBySub[sub] = subIDs
	return m.returnSubID, nil
}

func (m *mockSubStorage) ListTelegramIDsBySubscription(ctx context.Context, sub model.Subscription) ([]int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.idsBySub == nil {
		return nil, nil
	}

	return m.idsBySub[sub], nil
}

func (m *mockSubStorage) ListSubscriptionsByTelegramID(ctx context.Context, tgID int64) ([]model.Subscription, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.idsBySub == nil {
		return nil, nil
	}

	var subs []model.Subscription
	for sub, ids := range m.idsBySub {
		if slices.Contains(ids, tgID) {
			subs = append(subs, sub)
		}
	}
	return subs, nil
}

func assertNowIs(t *testing.T, expected time.Time) {
	t.Helper()
	assert.WithinDuration(t, expected, time.Now(), 10*time.Millisecond)
}
