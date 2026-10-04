package list

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
	ListBySeason(
		ctx context.Context,
		seasonID entity.SeasonID,
		memberID *entity.FamilyMemberID,
		status *entity.TaskCompletionStatus,
	) ([]entity.TaskCompletionView, error)
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

func (u *UseCase) List(ctx context.Context, req entity.TaskCompletionListRequest) ([]entity.TaskCompletionView, error) {
	member, err := u.familyMemberRepository.GetByUser(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("can't get family member: %w", err)
	}

	season, err := u.seasonRepository.GetByID(ctx, req.SeasonID)
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

	completions, err := u.taskCompletionRepository.ListBySeason(ctx, req.SeasonID, memberID, req.Status)
	if err != nil {
		return nil, fmt.Errorf("can't list task completions: %w", err)
	}

	return completions, nil
}
