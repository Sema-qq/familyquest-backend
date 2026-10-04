package create

import (
	"context"
	"fmt"

	"familyquest-backend/internal/domain"
	"familyquest-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type SeasonRepository interface {
	Create(ctx context.Context, season entity.SeasonCreate) error
}

type FamilyRepository interface {
	GetByUser(ctx context.Context, userID entity.UserID) (entity.FamilyWithRole, error)
}

type FamilyMemberRepository interface {
	GetByUser(ctx context.Context, userID entity.UserID) (entity.FamilyMember, error)
}

type UUIDGenerator interface {
	GenerateV7() uuid.UUID
}

type UseCase struct {
	seasonRepository       SeasonRepository
	familyRepository       FamilyRepository
	familyMemberRepository FamilyMemberRepository
	uuidGenerator          UUIDGenerator
}

func New(
	seasonRepository SeasonRepository,
	familyRepository FamilyRepository,
	familyMemberRepository FamilyMemberRepository,
	uuidGenerator UUIDGenerator,
) *UseCase {
	return &UseCase{
		seasonRepository:       seasonRepository,
		familyRepository:       familyRepository,
		familyMemberRepository: familyMemberRepository,
		uuidGenerator:          uuidGenerator,
	}
}

func (u *UseCase) Create(ctx context.Context, req entity.SeasonCreateRequest) (entity.Season, error) {
	member, err := u.familyMemberRepository.GetByUser(ctx, req.UserID)
	if err != nil {
		return entity.Season{}, fmt.Errorf("can't get family member: %w", err)
	}

	if member.Role != entity.FamilyRoleParent {
		return entity.Season{}, domain.ForbiddenError()
	}

	family, err := u.familyRepository.GetByUser(ctx, req.UserID)
	if err != nil {
		return entity.Season{}, fmt.Errorf("can't get family: %w", err)
	}

	season := entity.Season{
		ID:       entity.SeasonID(u.uuidGenerator.GenerateV7()),
		FamilyID: family.Family.ID,
		Title:    req.Title,
		StartsAt: req.StartsAt,
		EndsAt:   req.EndsAt,
		Timezone: family.Family.Timezone,
		Status:   entity.SeasonStatusDraft,
	}

	if err = u.seasonRepository.Create(ctx, entity.SeasonCreate{
		ID:       season.ID,
		FamilyID: season.FamilyID,
		Title:    season.Title,
		StartsAt: season.StartsAt,
		EndsAt:   season.EndsAt,
		Timezone: season.Timezone,
	}); err != nil {
		return entity.Season{}, fmt.Errorf("can't create season: %w", err)
	}

	return season, nil
}
