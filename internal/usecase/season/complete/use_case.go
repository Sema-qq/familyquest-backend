package complete

import (
	"context"
	"fmt"
	"time"

	"familyquest-backend/internal/domain"
	"familyquest-backend/internal/domain/entity"

	"github.com/guregu/null/v6"
)

type FamilyMemberRepository interface {
	GetByUser(ctx context.Context, userID entity.UserID) (entity.FamilyMember, error)
}

type SeasonRepository interface {
	GetByID(ctx context.Context, id entity.SeasonID) (entity.Season, error)
	Update(ctx context.Context, season entity.SeasonUpdate) error
}

type TaskCompletionRepository interface {
	CountSubmittedBySeason(ctx context.Context, seasonID entity.SeasonID) (int64, error)
}

type TimeProvider interface {
	Now() time.Time
}

type UseCase struct {
	familyMemberRepository   FamilyMemberRepository
	seasonRepository         SeasonRepository
	taskCompletionRepository TaskCompletionRepository
	timeProvider             TimeProvider
}

func New(
	familyMemberRepository FamilyMemberRepository,
	seasonRepository SeasonRepository,
	taskCompletionRepository TaskCompletionRepository,
	timeProvider TimeProvider,
) *UseCase {
	return &UseCase{
		familyMemberRepository:   familyMemberRepository,
		seasonRepository:         seasonRepository,
		taskCompletionRepository: taskCompletionRepository,
		timeProvider:             timeProvider,
	}
}

func (u *UseCase) Complete(ctx context.Context, userID entity.UserID, seasonID entity.SeasonID) (entity.Season, error) {
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

	if season.Status == entity.SeasonStatusCompleted {
		return season, nil
	}
	if season.Status != entity.SeasonStatusActive {
		return entity.Season{}, domain.Conflict("season must be active")
	}

	currentDate, err := u.currentDate(season.Timezone)
	if err != nil {
		return entity.Season{}, err
	}
	if !currentDate.After(season.EndsAt) {
		return entity.Season{}, domain.Conflict("season can be completed only after ends_at")
	}

	submittedCount, err := u.taskCompletionRepository.CountSubmittedBySeason(ctx, season.ID)
	if err != nil {
		return entity.Season{}, fmt.Errorf("can't count submitted task completions: %w", err)
	}
	if submittedCount > 0 {
		return entity.Season{}, domain.Conflict("season has submitted completions")
	}

	now := u.timeProvider.Now()
	season.Status = entity.SeasonStatusCompleted
	season.Completed = true
	if err = u.seasonRepository.Update(ctx, entity.SeasonUpdate{
		ID:          season.ID,
		Status:      null.ValueFrom(entity.SeasonStatusCompleted),
		CompletedAt: null.ValueFrom(now),
	}); err != nil {
		return entity.Season{}, fmt.Errorf("can't update season: %w", err)
	}

	return season, nil
}

func (u *UseCase) currentDate(timezone string) (time.Time, error) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("can't load season timezone: %w", err)
	}

	now := u.timeProvider.Now().In(location)
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC), nil
}
