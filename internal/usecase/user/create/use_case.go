package create

import (
	"context"
	"fmt"

	"familyquest-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user entity.UserCreate) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UUIDGenerator interface {
	GenerateV7() uuid.UUID
}

type UseCase struct {
	userRepository UserRepository
	passwordHasher PasswordHasher
	uuidGenerator  UUIDGenerator
}

func New(
	userRepository UserRepository,
	passwordHasher PasswordHasher,
	uuidGenerator UUIDGenerator,
) *UseCase {
	return &UseCase{
		userRepository: userRepository,
		passwordHasher: passwordHasher,
		uuidGenerator:  uuidGenerator,
	}
}

func (u *UseCase) Create(ctx context.Context, req entity.UserCreateRequest) (entity.User, error) {
	passwordHash, err := u.passwordHasher.Hash(req.Password)
	if err != nil {
		return entity.User{}, fmt.Errorf("hash password: %w", err)
	}

	user := entity.User{
		ID:           entity.UserID(u.uuidGenerator.GenerateV7()),
		Login:        req.Login,
		PasswordHash: passwordHash,
		DisplayName:  req.DisplayName,
	}

	if err = u.userRepository.Create(ctx, entity.UserCreate(user)); err != nil {
		return entity.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}
