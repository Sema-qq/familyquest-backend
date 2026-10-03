package membercreate

import (
	"context"
	"fmt"

	"familyquest-backend/internal/domain"
	"familyquest-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user entity.UserCreate) error
}

type FamilyMemberRepository interface {
	Create(ctx context.Context, member entity.FamilyMemberCreate) error
	GetByUser(ctx context.Context, userID entity.UserID) (entity.FamilyMember, error)
}

type TxManager interface {
	WithinTx(ctx context.Context, fn func(context.Context) error) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UUIDGenerator interface {
	GenerateV7() uuid.UUID
}

type UseCase struct {
	userRepository         UserRepository
	familyMemberRepository FamilyMemberRepository
	txManager              TxManager
	passwordHasher         PasswordHasher
	uuidGenerator          UUIDGenerator
}

func New(
	userRepository UserRepository,
	familyMemberRepository FamilyMemberRepository,
	txManager TxManager,
	passwordHasher PasswordHasher,
	uuidGenerator UUIDGenerator,
) *UseCase {
	return &UseCase{
		userRepository:         userRepository,
		familyMemberRepository: familyMemberRepository,
		txManager:              txManager,
		passwordHasher:         passwordHasher,
		uuidGenerator:          uuidGenerator,
	}
}

func (u *UseCase) Create(ctx context.Context, req entity.FamilyMemberCreateRequest) (entity.FamilyMemberProfile, error) {
	parentMember, err := u.familyMemberRepository.GetByUser(ctx, req.ParentID)
	if err != nil {
		return entity.FamilyMemberProfile{}, fmt.Errorf("can't get family member: %w", err)
	}

	if parentMember.Role != entity.FamilyRoleParent {
		return entity.FamilyMemberProfile{}, domain.ForbiddenError()
	}

	passwordHash, err := u.passwordHasher.Hash(req.Password)
	if err != nil {
		return entity.FamilyMemberProfile{}, fmt.Errorf("can't hash password: %w", err)
	}

	user := entity.User{
		ID:          entity.UserID(u.uuidGenerator.GenerateV7()),
		Login:       req.Login,
		DisplayName: req.DisplayName,
	}
	member := entity.FamilyMember{
		ID:       entity.FamilyMemberID(u.uuidGenerator.GenerateV7()),
		FamilyID: parentMember.FamilyID,
		UserID:   user.ID,
		Role:     entity.FamilyRoleChild,
	}

	if err = u.txManager.WithinTx(ctx, func(ctx context.Context) error {
		if err := u.userRepository.Create(ctx, entity.UserCreate{
			ID:           user.ID,
			Login:        user.Login,
			PasswordHash: passwordHash,
			DisplayName:  user.DisplayName,
		}); err != nil {
			return fmt.Errorf("can't create user: %w", err)
		}

		if err := u.familyMemberRepository.Create(ctx, entity.FamilyMemberCreate{
			ID:       member.ID,
			FamilyID: member.FamilyID,
			UserID:   member.UserID,
			Role:     member.Role,
		}); err != nil {
			return fmt.Errorf("can't create family member: %w", err)
		}

		return nil
	}); err != nil {
		return entity.FamilyMemberProfile{}, fmt.Errorf("can't create family member: %w", err)
	}

	return entity.FamilyMemberProfile{
		ID:          member.ID,
		FamilyID:    member.FamilyID,
		UserID:      user.ID,
		Login:       user.Login,
		DisplayName: user.DisplayName,
		Role:        member.Role,
	}, nil
}
