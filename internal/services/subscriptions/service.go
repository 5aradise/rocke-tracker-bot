package subservice

import (
	"bot/config"
	model "bot/internal/models"
	"context"
	"errors"
	"fmt"
)

type Service struct {
	api  rocketLeagueAPI
	subs subStorage
}

type rocketLeagueAPI interface {
	Tournaments() ([]model.Tournament, error)
}

type subStorage interface {
	CreateSubscriptionByTelegramID(ctx context.Context,
		tgID int64, sub model.Subscription) (int64, error)
	ListSubscriptionsByTelegramID(ctx context.Context,
		tgID int64) ([]model.Subscription, error)
	ListTelegramIDsBySubscription(ctx context.Context,
		sub model.Subscription) ([]int64, error)
}

func New(api rocketLeagueAPI, subStor subStorage) *Service {
	return &Service{
		api:  api,
		subs: subStor,
	}
}

// Codes: [config.CodeUserHasSub], [config.CodeUserWithTgIDNotExist]
func (s *Service) SubscribeByTelegram(ctx context.Context,
	tgID int64, sub model.Subscription) (int64, config.Error) {
	id, err := s.subs.CreateSubscriptionByTelegramID(ctx, tgID, sub)
	if err != nil {
		if errors.Is(err, config.ErrUniqueConstraint) {
			return 0, config.NewError(
				config.CodeUserHasSub,
				fmt.Errorf("%w: user already has this subscription", err),
			)
		}
		if errors.Is(err, config.ErrNotFound) {
			return 0, config.NewError(
				config.CodeUserWithTgIDNotExist,
				fmt.Errorf("%w: user with this Telegram ID does not exist", err),
			)
		}
		return 0, config.NewUnknownError(err)
	}

	return id, config.NilError
}

func (s *Service) ListTelegramUserSubscriptions(ctx context.Context,
	tgID int64) ([]model.Subscription, config.Error) {
	subs, err := s.subs.ListSubscriptionsByTelegramID(ctx, tgID)
	if err != nil {
		return nil, config.NewUnknownError(err)
	}

	return subs, config.NilError
}
