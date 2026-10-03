package get

import (
	"context"
	"fmt"

	"familyquest-backend/internal/domain/entity"
)

type UserRepository interface {
	GetByID(ctx context.Context, id entity.UserID) (entity.User, error)
}

type UseCase struct {
	userRepository UserRepository
}

func New(userRepository UserRepository) *UseCase {
	return &UseCase{
		userRepository: userRepository,
	}
}

func (u *UseCase) Get(ctx context.Context, id entity.UserID) (entity.User, error) {
	user, err := u.userRepository.GetByID(ctx, id)
	if err != nil {
		return entity.User{}, fmt.Errorf("can't get user: %w", err)
	}

	return user, nil
}
