package taskadd

import (
	"context"
	"fmt"

	"familyquest-backend/internal/domain"
	"familyquest-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type FamilyMemberRepository interface {
	GetByUser(ctx context.Context, userID entity.UserID) (entity.FamilyMember, error)
}

type SeasonRepository interface {
	GetByID(ctx context.Context, id entity.SeasonID) (entity.Season, error)
}

type TaskRepository interface {
	GetByID(ctx context.Context, id entity.TaskID) (entity.Task, error)
}

type SeasonTaskRepository interface {
	Create(ctx context.Context, task entity.SeasonTaskCreate) error
}

type SeasonTaskSlotRepository interface {
	CreateMany(ctx context.Context, slots []entity.SeasonTaskSlotCreate) error
}

type TxManager interface {
	WithinTx(ctx context.Context, fn func(context.Context) error) error
}

type UUIDGenerator interface {
	GenerateV7() uuid.UUID
}

type UseCase struct {
	familyMemberRepository   FamilyMemberRepository
	seasonRepository         SeasonRepository
	taskRepository           TaskRepository
	seasonTaskRepository     SeasonTaskRepository
	seasonTaskSlotRepository SeasonTaskSlotRepository
	txManager                TxManager
	uuidGenerator            UUIDGenerator
}

func New(
	familyMemberRepository FamilyMemberRepository,
	seasonRepository SeasonRepository,
	taskRepository TaskRepository,
	seasonTaskRepository SeasonTaskRepository,
	seasonTaskSlotRepository SeasonTaskSlotRepository,
	txManager TxManager,
	uuidGenerator UUIDGenerator,
) *UseCase {
	return &UseCase{
		familyMemberRepository:   familyMemberRepository,
		seasonRepository:         seasonRepository,
		taskRepository:           taskRepository,
		seasonTaskRepository:     seasonTaskRepository,
		seasonTaskSlotRepository: seasonTaskSlotRepository,
		txManager:                txManager,
		uuidGenerator:            uuidGenerator,
	}
}

func (u *UseCase) Add(ctx context.Context, req entity.SeasonTaskAddRequest) (entity.SeasonTask, error) {
	var result entity.SeasonTask

	if err := u.txManager.WithinTx(ctx, func(ctx context.Context) error {
		member, err := u.familyMemberRepository.GetByUser(ctx, req.UserID)
		if err != nil {
			return fmt.Errorf("can't get family member: %w", err)
		}

		if member.Role != entity.FamilyRoleParent {
			return domain.ForbiddenError()
		}

		season, err := u.seasonRepository.GetByID(ctx, req.SeasonID)
		if err != nil {
			return fmt.Errorf("can't get season: %w", err)
		}

		if season.FamilyID != member.FamilyID {
			return domain.NotFound("season not found")
		}
		if season.Status != entity.SeasonStatusDraft {
			return domain.Conflict("season must be draft")
		}

		task, err := u.taskRepository.GetByID(ctx, req.TaskID)
		if err != nil {
			return fmt.Errorf("can't get task: %w", err)
		}

		if task.FamilyID != season.FamilyID {
			return domain.NotFound("task not found")
		}

		slots, err := u.buildSlots(req.Schedule, season)
		if err != nil {
			return err
		}

		seasonTask := entity.SeasonTask{
			ID:                entity.SeasonTaskID(u.uuidGenerator.GenerateV7()),
			SeasonID:          season.ID,
			TaskID:            task.ID,
			Title:             task.Title,
			Description:       task.Description,
			Points:            task.Points,
			ParticipationMode: req.ParticipationMode,
			Slots:             slots,
		}

		if err = u.seasonTaskRepository.Create(ctx, entity.SeasonTaskCreate{
			ID:                seasonTask.ID,
			SeasonID:          seasonTask.SeasonID,
			TaskID:            seasonTask.TaskID,
			Title:             seasonTask.Title,
			Description:       seasonTask.Description,
			Points:            seasonTask.Points,
			ParticipationMode: seasonTask.ParticipationMode,
		}); err != nil {
			return fmt.Errorf("can't create season task: %w", err)
		}

		slotCreates := make([]entity.SeasonTaskSlotCreate, len(seasonTask.Slots))
		for i, slot := range seasonTask.Slots {
			seasonTask.Slots[i].SeasonTaskID = seasonTask.ID
			slotCreates[i] = entity.SeasonTaskSlotCreate{
				ID:             slot.ID,
				SeasonTaskID:   seasonTask.ID,
				Position:       slot.Position,
				AvailableFrom:  slot.AvailableFrom,
				AvailableUntil: slot.AvailableUntil,
			}
		}

		if err = u.seasonTaskSlotRepository.CreateMany(ctx, slotCreates); err != nil {
			return fmt.Errorf("can't create season task slots: %w", err)
		}

		result = seasonTask
		return nil
	}); err != nil {
		return entity.SeasonTask{}, fmt.Errorf("can't add task to season: %w", err)
	}

	return result, nil
}

func (u *UseCase) buildSlots(schedule entity.SeasonTaskSchedule, season entity.Season) ([]entity.SeasonTaskSlot, error) {
	switch schedule.Type {
	case entity.SeasonTaskScheduleTypeDates:
		slots := make([]entity.SeasonTaskSlot, len(schedule.Dates))
		for i, date := range schedule.Dates {
			if date.Before(season.StartsAt) || date.After(season.EndsAt) {
				return nil, domain.ValidationError("schedule dates must be inside season")
			}

			slots[i] = entity.SeasonTaskSlot{
				ID:             entity.SeasonTaskSlotID(u.uuidGenerator.GenerateV7()),
				Position:       int64(i + 1),
				AvailableFrom:  date,
				AvailableUntil: date,
			}
		}

		return slots, nil
	case entity.SeasonTaskScheduleTypeQuota:
		slots := make([]entity.SeasonTaskSlot, schedule.Count)
		for i := range slots {
			slots[i] = entity.SeasonTaskSlot{
				ID:             entity.SeasonTaskSlotID(u.uuidGenerator.GenerateV7()),
				Position:       int64(i + 1),
				AvailableFrom:  season.StartsAt,
				AvailableUntil: season.EndsAt,
			}
		}

		return slots, nil
	default:
		return nil, domain.ValidationError("unsupported schedule type")
	}
}
