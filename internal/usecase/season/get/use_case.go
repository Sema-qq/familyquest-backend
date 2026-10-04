package get

import (
	"context"
	"fmt"

	"familyquest-backend/internal/domain"
	"familyquest-backend/internal/domain/entity"
)

type SeasonRepository interface {
	GetByID(ctx context.Context, id entity.SeasonID) (entity.Season, error)
}

type SeasonTaskRepository interface {
	ListBySeason(ctx context.Context, seasonID entity.SeasonID) ([]entity.SeasonTask, error)
}

type SeasonTaskSlotRepository interface {
	ListBySeason(ctx context.Context, seasonID entity.SeasonID) ([]entity.SeasonTaskSlot, error)
}

type TaskCompletionRepository interface {
	ListBySeason(
		ctx context.Context,
		seasonID entity.SeasonID,
		memberID *entity.FamilyMemberID,
		status *entity.TaskCompletionStatus,
	) ([]entity.TaskCompletionView, error)
}

type FamilyMemberRepository interface {
	GetByUser(ctx context.Context, userID entity.UserID) (entity.FamilyMember, error)
}

type UseCase struct {
	seasonRepository         SeasonRepository
	familyMemberRepository   FamilyMemberRepository
	seasonTaskRepository     SeasonTaskRepository
	seasonTaskSlotRepository SeasonTaskSlotRepository
	taskCompletionRepository TaskCompletionRepository
}

func New(
	seasonRepository SeasonRepository,
	familyMemberRepository FamilyMemberRepository,
	seasonTaskRepository SeasonTaskRepository,
	seasonTaskSlotRepository SeasonTaskSlotRepository,
	taskCompletionRepository TaskCompletionRepository,
) *UseCase {
	return &UseCase{
		seasonRepository:         seasonRepository,
		familyMemberRepository:   familyMemberRepository,
		seasonTaskRepository:     seasonTaskRepository,
		seasonTaskSlotRepository: seasonTaskSlotRepository,
		taskCompletionRepository: taskCompletionRepository,
	}
}

func (u *UseCase) Get(ctx context.Context, userID entity.UserID, seasonID entity.SeasonID) (entity.SeasonDetail, error) {
	member, err := u.familyMemberRepository.GetByUser(ctx, userID)
	if err != nil {
		return entity.SeasonDetail{}, fmt.Errorf("can't get family member: %w", err)
	}

	season, err := u.seasonRepository.GetByID(ctx, seasonID)
	if err != nil {
		return entity.SeasonDetail{}, fmt.Errorf("can't get season: %w", err)
	}

	if season.FamilyID != member.FamilyID {
		return entity.SeasonDetail{}, domain.NotFound("season not found")
	}

	tasks, err := u.seasonTaskRepository.ListBySeason(ctx, season.ID)
	if err != nil {
		return entity.SeasonDetail{}, fmt.Errorf("can't list season tasks: %w", err)
	}

	slots, err := u.seasonTaskSlotRepository.ListBySeason(ctx, season.ID)
	if err != nil {
		return entity.SeasonDetail{}, fmt.Errorf("can't list season task slots: %w", err)
	}

	var memberID *entity.FamilyMemberID
	if member.Role == entity.FamilyRoleChild {
		memberID = &member.ID
	}
	completions, err := u.taskCompletionRepository.ListBySeason(ctx, season.ID, memberID, nil)
	if err != nil {
		return entity.SeasonDetail{}, fmt.Errorf("can't list task completions: %w", err)
	}

	return buildDetail(season, tasks, slots, completions), nil
}

func buildDetail(
	season entity.Season,
	tasks []entity.SeasonTask,
	slots []entity.SeasonTaskSlot,
	completions []entity.TaskCompletionView,
) entity.SeasonDetail {
	completionsBySlot := make(map[entity.SeasonTaskSlotID][]entity.TaskCompletion)
	for _, completion := range completions {
		completionsBySlot[completion.SlotID] = append(completionsBySlot[completion.SlotID], entity.TaskCompletion{
			ID:            completion.ID,
			MemberID:      completion.MemberID,
			SlotID:        completion.SlotID,
			SeasonTaskID:  completion.SeasonTaskID,
			PerformedOn:   completion.PerformedOn,
			Status:        completion.Status,
			ChildComment:  completion.ChildComment,
			ParentComment: completion.ParentComment,
			ReviewedBy:    completion.ReviewedBy,
			ReviewedAt:    completion.ReviewedAt,
		})
	}

	slotsByTask := make(map[entity.SeasonTaskID][]entity.SeasonTaskSlotDetail)
	for _, slot := range slots {
		slotsByTask[slot.SeasonTaskID] = append(slotsByTask[slot.SeasonTaskID], entity.SeasonTaskSlotDetail{
			SeasonTaskSlot: slot,
			Completions:    completionsBySlot[slot.ID],
		})
	}

	result := entity.SeasonDetail{
		Season: season,
		Tasks:  make([]entity.SeasonTaskDetail, len(tasks)),
	}
	for i, task := range tasks {
		result.Tasks[i] = entity.SeasonTaskDetail{
			SeasonTask: task,
			Slots:      slotsByTask[task.ID],
		}
	}

	return result
}
