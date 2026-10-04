package activate

import (
	"context"
	"fmt"

	"familyquest-backend/internal/domain"
	"familyquest-backend/internal/domain/entity"

	"github.com/guregu/null/v6"
)

type SeasonRepository interface {
	GetByID(ctx context.Context, id entity.SeasonID) (entity.Season, error)
	Update(ctx context.Context, season entity.SeasonUpdate) error
}

type FamilyMemberRepository interface {
	GetByUser(ctx context.Context, userID entity.UserID) (entity.FamilyMember, error)
}

type SeasonTaskRepository interface {
	CountBySeason(ctx context.Context, seasonID entity.SeasonID) (int64, error)
}

type SeasonTaskSlotRepository interface {
	CountBySeason(ctx context.Context, seasonID entity.SeasonID) (int64, error)
}

type UseCase struct {
	seasonRepository         SeasonRepository
	familyMemberRepository   FamilyMemberRepository
	seasonTaskRepository     SeasonTaskRepository
	seasonTaskSlotRepository SeasonTaskSlotRepository
}

func New(
	seasonRepository SeasonRepository,
	familyMemberRepository FamilyMemberRepository,
	seasonTaskRepository SeasonTaskRepository,
	seasonTaskSlotRepository SeasonTaskSlotRepository,
) *UseCase {
	return &UseCase{
		seasonRepository:         seasonRepository,
		familyMemberRepository:   familyMemberRepository,
		seasonTaskRepository:     seasonTaskRepository,
		seasonTaskSlotRepository: seasonTaskSlotRepository,
	}
}

func (u *UseCase) Activate(ctx context.Context, userID entity.UserID, seasonID entity.SeasonID) (entity.Season, error) {
	member, err := u.familyMemberRepository.GetByUser(ctx, userID)
	if err != nil {
		return entity.Season{}, fmt.Errorf("can't get family member: %w", err)
	}

	if member.Role != entity.FamilyRoleParent {
		return entity.Season{}, domain.ForbiddenError()
	}

	season, err := u.seasonRepository.GetByID(ctx, seasonID)
	if err != nil {
		return entity.Season{}, fmt.Errorf("can't get season: %w", err)
	}

	if season.FamilyID != member.FamilyID {
		return entity.Season{}, domain.NotFound("season not found")
	}

	switch season.Status {
	case entity.SeasonStatusActive:
		return season, nil
	case entity.SeasonStatusCompleted:
		return entity.Season{}, domain.Conflict("season already completed")
	}

	taskCount, err := u.seasonTaskRepository.CountBySeason(ctx, season.ID)
	if err != nil {
		return entity.Season{}, fmt.Errorf("can't count season tasks: %w", err)
	}
	if taskCount == 0 {
		return entity.Season{}, domain.Conflict("season must contain tasks")
	}

	slotCount, err := u.seasonTaskSlotRepository.CountBySeason(ctx, season.ID)
	if err != nil {
		return entity.Season{}, fmt.Errorf("can't count season task slots: %w", err)
	}
	if slotCount == 0 {
		return entity.Season{}, domain.Conflict("season must contain task slots")
	}

	season.Status = entity.SeasonStatusActive
	if err = u.seasonRepository.Update(ctx, entity.SeasonUpdate{
		ID:     season.ID,
		Status: null.ValueFrom(entity.SeasonStatusActive),
	}); err != nil {
		return entity.Season{}, fmt.Errorf("can't update season: %w", err)
	}

	return season, nil
}
