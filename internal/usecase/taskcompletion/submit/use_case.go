package submit

import (
	"context"
	"fmt"
	"time"

	"familyquest-backend/internal/domain"
	"familyquest-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type FamilyMemberRepository interface {
	GetByUser(ctx context.Context, userID entity.UserID) (entity.FamilyMember, error)
}

type SeasonTaskRepository interface {
	GetByID(ctx context.Context, id entity.SeasonTaskID) (entity.SeasonTaskWithSeason, error)
}

type SeasonTaskSlotRepository interface {
	FindAvailable(
		ctx context.Context,
		seasonTaskID entity.SeasonTaskID,
		memberID entity.FamilyMemberID,
		performedOn time.Time,
		shared bool,
	) (entity.SeasonTaskSlot, error)
}

type TaskCompletionRepository interface {
	Create(ctx context.Context, completion entity.TaskCompletionCreate) (entity.TaskCompletion, error)
	GetRejectedByMemberAndSlot(
		ctx context.Context,
		memberID entity.FamilyMemberID,
		slotID entity.SeasonTaskSlotID,
	) (entity.TaskCompletion, bool, error)
	ResubmitRejected(ctx context.Context, completion entity.TaskCompletionResubmit) (entity.TaskCompletion, error)
}

type TxManager interface {
	WithinTx(ctx context.Context, fn func(context.Context) error) error
}

type TimeProvider interface {
	Now() time.Time
}

type UUIDGenerator interface {
	GenerateV7() uuid.UUID
}

type UseCase struct {
	familyMemberRepository   FamilyMemberRepository
	seasonTaskRepository     SeasonTaskRepository
	seasonTaskSlotRepository SeasonTaskSlotRepository
	taskCompletionRepository TaskCompletionRepository
	txManager                TxManager
	timeProvider             TimeProvider
	uuidGenerator            UUIDGenerator
}

func New(
	familyMemberRepository FamilyMemberRepository,
	seasonTaskRepository SeasonTaskRepository,
	seasonTaskSlotRepository SeasonTaskSlotRepository,
	taskCompletionRepository TaskCompletionRepository,
	txManager TxManager,
	timeProvider TimeProvider,
	uuidGenerator UUIDGenerator,
) *UseCase {
	return &UseCase{
		familyMemberRepository:   familyMemberRepository,
		seasonTaskRepository:     seasonTaskRepository,
		seasonTaskSlotRepository: seasonTaskSlotRepository,
		taskCompletionRepository: taskCompletionRepository,
		txManager:                txManager,
		timeProvider:             timeProvider,
		uuidGenerator:            uuidGenerator,
	}
}

func (u *UseCase) Submit(ctx context.Context, req entity.TaskCompletionSubmitRequest) (entity.TaskCompletion, error) {
	var result entity.TaskCompletion

	if err := u.txManager.WithinTx(ctx, func(ctx context.Context) error {
		member, err := u.familyMemberRepository.GetByUser(ctx, req.UserID)
		if err != nil {
			return fmt.Errorf("can't get family member: %w", err)
		}

		if member.Role != entity.FamilyRoleChild {
			return domain.ForbiddenError()
		}

		seasonTask, err := u.seasonTaskRepository.GetByID(ctx, req.SeasonTaskID)
		if err != nil {
			return fmt.Errorf("can't get season task: %w", err)
		}

		if seasonTask.SeasonFamilyID != member.FamilyID {
			return domain.NotFound("season task not found")
		}
		if seasonTask.SeasonStatus != entity.SeasonStatusActive {
			return domain.Conflict("season must be active")
		}

		currentDate, err := u.currentDate(seasonTask.SeasonTimezone)
		if err != nil {
			return err
		}

		if currentDate.Before(seasonTask.SeasonStartsAt) || currentDate.After(seasonTask.SeasonEndsAt) {
			return domain.Conflict("current date must be inside season")
		}
		if req.PerformedOn.Before(seasonTask.SeasonStartsAt) || req.PerformedOn.After(seasonTask.SeasonEndsAt) {
			return domain.ValidationError(`the "performed_on" field must be inside season`)
		}
		if req.PerformedOn.After(currentDate) {
			return domain.ValidationError(`the "performed_on" field must not be after current date`)
		}

		slot, err := u.seasonTaskSlotRepository.FindAvailable(
			ctx,
			seasonTask.ID,
			member.ID,
			req.PerformedOn,
			seasonTask.ParticipationMode == entity.TaskParticipationModeShared,
		)
		if err != nil {
			return fmt.Errorf("can't find available season task slot: %w", err)
		}

		rejected, ok, err := u.taskCompletionRepository.GetRejectedByMemberAndSlot(ctx, member.ID, slot.ID)
		if err != nil {
			return fmt.Errorf("can't get rejected task completion: %w", err)
		}
		if ok {
			result, err = u.taskCompletionRepository.ResubmitRejected(ctx, entity.TaskCompletionResubmit{
				ID:           rejected.ID,
				PerformedOn:  req.PerformedOn,
				ChildComment: req.Comment,
			})
			if err != nil {
				return fmt.Errorf("can't resubmit rejected task completion: %w", err)
			}

			return nil
		}

		result, err = u.taskCompletionRepository.Create(ctx, entity.TaskCompletionCreate{
			ID:           entity.TaskCompletionID(u.uuidGenerator.GenerateV7()),
			MemberID:     member.ID,
			SlotID:       slot.ID,
			PerformedOn:  req.PerformedOn,
			ChildComment: req.Comment,
		})
		if err != nil {
			return fmt.Errorf("can't create task completion: %w", err)
		}

		return nil
	}); err != nil {
		return entity.TaskCompletion{}, fmt.Errorf("can't submit task completion: %w", err)
	}

	return result, nil
}

func (u *UseCase) currentDate(timezone string) (time.Time, error) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("can't load season timezone: %w", err)
	}

	now := u.timeProvider.Now().In(location)
	return dateOnly(now), nil
}

func dateOnly(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC)
}
