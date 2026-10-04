package approve

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
	Approve(ctx context.Context, update entity.TaskCompletionUpdateReview) (entity.TaskCompletion, error)
	HasApprovedBySlot(ctx context.Context, slotID entity.SeasonTaskSlotID) (bool, error)
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

func (u *UseCase) Approve(ctx context.Context, req entity.TaskCompletionReviewRequest) (entity.TaskCompletion, error) {
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
	case entity.TaskCompletionStatusApproved:
		return completion.TaskCompletion, nil
	case entity.TaskCompletionStatusRejected:
		return entity.TaskCompletion{}, domain.Conflict("rejected completion must be submitted again")
	}

	if completion.ParticipationMode == entity.TaskParticipationModeShared {
		hasApproved, err := u.taskCompletionRepository.HasApprovedBySlot(ctx, completion.SlotID)
		if err != nil {
			return entity.TaskCompletion{}, fmt.Errorf("can't check approved task completion: %w", err)
		}
		if hasApproved {
			return entity.TaskCompletion{}, domain.Conflict("slot already approved")
		}
	}

	result, err := u.taskCompletionRepository.Approve(ctx, entity.TaskCompletionUpdateReview{
		ID:         completion.ID,
		Status:     entity.TaskCompletionStatusApproved,
		ReviewedBy: parent.ID,
		ReviewedAt: u.timeProvider.Now(),
	})
	if err != nil {
		return entity.TaskCompletion{}, fmt.Errorf("can't approve task completion: %w", err)
	}

	return result, nil
}
