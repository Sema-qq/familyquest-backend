package results

import (
	"context"
	"fmt"

	"familyquest-backend/internal/domain"
	"familyquest-backend/internal/domain/entity"
)

type FamilyMemberRepository interface {
	GetByUser(ctx context.Context, userID entity.UserID) (entity.FamilyMember, error)
}

type SeasonRepository interface {
	GetByID(ctx context.Context, id entity.SeasonID) (entity.Season, error)
}

type TaskCompletionRepository interface {
	ResultsBySeason(
		ctx context.Context,
		seasonID entity.SeasonID,
		memberID *entity.FamilyMemberID,
	) ([]entity.SeasonResult, error)
}

type UseCase struct {
	familyMemberRepository   FamilyMemberRepository
	seasonRepository         SeasonRepository
	taskCompletionRepository TaskCompletionRepository
}

func New(
	familyMemberRepository FamilyMemberRepository,
	seasonRepository SeasonRepository,
	taskCompletionRepository TaskCompletionRepository,
) *UseCase {
	return &UseCase{
		familyMemberRepository:   familyMemberRepository,
		seasonRepository:         seasonRepository,
		taskCompletionRepository: taskCompletionRepository,
	}
}

func (u *UseCase) Get(ctx context.Context, userID entity.UserID, seasonID entity.SeasonID) ([]entity.SeasonResult, error) {
	member, err := u.familyMemberRepository.GetByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("can't get family member: %w", err)
	}

	season, err := u.seasonRepository.GetByID(ctx, seasonID)
	if err != nil {
		return nil, fmt.Errorf("can't get season: %w", err)
	}
	if season.FamilyID != member.FamilyID {
		return nil, domain.NotFound("season not found")
	}

	var memberID *entity.FamilyMemberID
	if member.Role == entity.FamilyRoleChild {
		memberID = &member.ID
	}

	results, err := u.taskCompletionRepository.ResultsBySeason(ctx, seasonID, memberID)
	if err != nil {
		return nil, fmt.Errorf("can't get season results: %w", err)
	}

	return results, nil
}
