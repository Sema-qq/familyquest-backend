package reject

import (
	"context"
	"fmt"
	"time"

	"familyquest-backend/internal/domain"
	"familyquest-backend/internal/domain/entity"
)

type FamilyMemberRepository interface {
	GetByUser(ctx context.Context, userID entity.UserID) (entity.FamilyMember, error)
}

type TaskCompletionRepository interface {
	GetReviewByID(ctx context.Context, id entity.TaskCompletionID) (entity.TaskCompletionReview, error)
	Reject(ctx context.Context, update entity.TaskCompletionUpdateReview) (entity.TaskCompletion, error)
}

type TimeProvider interface {
	Now() time.Time
}

type UseCase struct {
	familyMemberRepository   FamilyMemberRepository
	taskCompletionRepository TaskCompletionRepository
	timeProvider             TimeProvider
}

func New(
	familyMemberRepository FamilyMemberRepository,
	taskCompletionRepository TaskCompletionRepository,
	timeProvider TimeProvider,
) *UseCase {
	return &UseCase{
		familyMemberRepository:   familyMemberRepository,
		taskCompletionRepository: taskCompletionRepository,
		timeProvider:             timeProvider,
	}
}

func (u *UseCase) Reject(ctx context.Context, req entity.TaskCompletionReviewRequest) (entity.TaskCompletion, error) {
	parent, err := u.familyMemberRepository.GetByUser(ctx, req.UserID)
	if err != nil {
		return entity.TaskCompletion{}, fmt.Errorf("can't get family member: %w", err)
	}
	if parent.Role != entity.FamilyRoleParent {
		return entity.TaskCompletion{}, domain.ForbiddenError()
	}

	completion, err := u.taskCompletionRepository.GetReviewByID(ctx, req.CompletionID)
	if err != nil {
		return entity.TaskCompletion{}, fmt.Errorf("can't get task completion: %w", err)
	}
	if completion.SeasonFamilyID != parent.FamilyID {
		return entity.TaskCompletion{}, domain.NotFound("task completion not found")
	}
	if completion.SeasonStatus != entity.SeasonStatusActive {
		return entity.TaskCompletion{}, domain.Conflict("season must be active")
	}

	switch completion.Status {
	case entity.TaskCompletionStatusRejected:
		return completion.TaskCompletion, nil
	case entity.TaskCompletionStatusApproved:
		return entity.TaskCompletion{}, domain.Conflict("approved completion can't be rejected")
	}

	result, err := u.taskCompletionRepository.Reject(ctx, entity.TaskCompletionUpdateReview{
		ID:            completion.ID,
		Status:        entity.TaskCompletionStatusRejected,
		ParentComment: req.Comment,
		ReviewedBy:    parent.ID,
		ReviewedAt:    u.timeProvider.Now(),
	})
	if err != nil {
		return entity.TaskCompletion{}, fmt.Errorf("can't reject task completion: %w", err)
	}

	return result, nil
}
