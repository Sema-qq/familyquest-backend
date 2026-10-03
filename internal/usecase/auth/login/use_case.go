package login

import (
	"context"
	"fmt"
	"time"

	"familyquest-backend/internal/domain"
	"familyquest-backend/internal/domain/entity"
)

type UserRepository interface {
	GetByLogin(ctx context.Context, login string) (entity.UserCredentials, error)
}

type PasswordHasher interface {
	Compare(hash string, password string) bool
}

type TimeProvider interface {
	Now() time.Time
}

type TokenGenerator interface {
	Generate(user entity.User, now time.Time) (string, time.Time, error)
}

type UseCase struct {
	userRepository UserRepository
	passwordHasher PasswordHasher
	timeProvider   TimeProvider
	tokenGenerator TokenGenerator
}

func New(
	userRepository UserRepository,
	passwordHasher PasswordHasher,
	timeProvider TimeProvider,
	tokenGenerator TokenGenerator,
) *UseCase {
	return &UseCase{
		userRepository: userRepository,
		passwordHasher: passwordHasher,
		timeProvider:   timeProvider,
		tokenGenerator: tokenGenerator,
	}
}

func (u *UseCase) Login(ctx context.Context, req entity.AuthRequest) (entity.AuthResult, error) {
	credentials, err := u.userRepository.GetByLogin(ctx, req.Login)
	if err != nil {
		return entity.AuthResult{}, fmt.Errorf("can't get user by login: %w", err)
	}

	if !u.passwordHasher.Compare(credentials.PasswordHash, req.Password) {
		return entity.AuthResult{}, domain.AuthorizationError()
	}

	accessToken, expiresAt, err := u.tokenGenerator.Generate(credentials.User, u.timeProvider.Now())
	if err != nil {
		return entity.AuthResult{}, fmt.Errorf("can't generate access token: %w", err)
	}

	return entity.AuthResult{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
		User:        credentials.User,
	}, nil
}
