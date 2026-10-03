package create

import (
	"context"
	"fmt"

	"familyquest-backend/internal/domain"
	"familyquest-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type FamilyRepository interface {
	Create(ctx context.Context, family entity.FamilyCreate) error
}

type FamilyMemberRepository interface {
	Create(ctx context.Context, member entity.FamilyMemberCreate) error
	GetByUser(ctx context.Context, userID entity.UserID) (entity.FamilyMember, error)
}

type TxManager interface {
	WithinTx(ctx context.Context, fn func(context.Context) error) error
}

type UUIDGenerator interface {
	GenerateV7() uuid.UUID
}

type UseCase struct {
	familyRepository       FamilyRepository
	familyMemberRepository FamilyMemberRepository
	txManager              TxManager
	uuidGenerator          UUIDGenerator
}

func New(
	familyRepository FamilyRepository,
	familyMemberRepository FamilyMemberRepository,
	txManager TxManager,
	uuidGenerator UUIDGenerator,
) *UseCase {
	return &UseCase{
		familyRepository:       familyRepository,
		familyMemberRepository: familyMemberRepository,
		txManager:              txManager,
		uuidGenerator:          uuidGenerator,
	}
}

func (u *UseCase) Create(ctx context.Context, req entity.FamilyCreateRequest) (entity.FamilyWithRole, error) {
	if _, err := u.familyMemberRepository.GetByUser(ctx, req.OwnerID); err == nil {
		return entity.FamilyWithRole{}, domain.Conflict("user already has family")
	}

	family := entity.Family{
		ID:       entity.FamilyID(u.uuidGenerator.GenerateV7()),
		Name:     req.Name,
		Timezone: req.Timezone,
		OwnerID:  req.OwnerID,
	}
	member := entity.FamilyMemberCreate{
		ID:       entity.FamilyMemberID(u.uuidGenerator.GenerateV7()),
		FamilyID: family.ID,
		UserID:   req.OwnerID,
		Role:     entity.FamilyRoleParent,
	}

	if err := u.txManager.WithinTx(ctx, func(ctx context.Context) error {
		if err := u.familyRepository.Create(ctx, entity.FamilyCreate(family)); err != nil {
			return fmt.Errorf("can't create family: %w", err)
		}

		if err := u.familyMemberRepository.Create(ctx, member); err != nil {
			return fmt.Errorf("can't create family member: %w", err)
		}

		return nil
	}); err != nil {
		return entity.FamilyWithRole{}, err
	}

	return entity.FamilyWithRole{
		Family: family,
		Role:   entity.FamilyRoleParent,
	}, nil
}
