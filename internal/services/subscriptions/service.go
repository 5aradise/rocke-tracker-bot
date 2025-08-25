package subService

import (
	"bot/config"
	model "bot/internal/models"
	subStorage "bot/internal/storage/subscriptions"
	"context"
	"errors"
	"fmt"
)

type Service struct {
	subs subStorage.Storage
}

func New(subStor subStorage.Storage) *Service {
	return &Service{subStor}
}

// Codes: [config.CodeUserHasSub], [config.CodeUserWithTgIDNotExist]
func (s *Service) SubscribeByTelegram(ctx context.Context, tgID int64, sub model.Subscription) (int64, config.Error) {
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
