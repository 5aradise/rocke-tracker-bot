package userService

import (
	"bot/config"
	userStorage "bot/internal/storage/users"
	"context"
	"errors"
	"fmt"
)

type Service struct {
	users userStorage.Storage
}

func New(userStor userStorage.Storage) *Service {
	return &Service{userStor}
}

// Codes: [config.CodeUserWithTgIDExist]
func (s *Service) CreateUser(ctx context.Context, tgID int64) (int64, config.Error) {
	id, err := s.users.CreateUser(ctx, tgID)
	if err != nil {
		if errors.Is(err, config.ErrUniqueConstraint) {
			return 0, config.NewError(
				config.CodeUserWithTgIDExist,
				fmt.Errorf("%w: user with this Telegram ID already exists", err),
			)
		}

		return 0, config.NewUnknownError(err)
	}

	return id, config.NilError
}
