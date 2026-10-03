package add

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"familyquest-backend/internal/domain"
	"familyquest-backend/internal/domain/entity"
)

type FamilyMemberRepository interface {
	GetByUser(ctx context.Context, userID entity.UserID) (entity.FamilyMember, error)
}

type TaskRepository interface {
	Create(ctx context.Context, item entity.TaskCreate) error
}

type UUIDGenerator interface {
	GenerateV7() uuid.UUID
}

type UseCase struct {
	taskRepository         TaskRepository
	familyMemberRepository FamilyMemberRepository
	uuidGenerator          UUIDGenerator
}

func New(taskRepository TaskRepository, familyMemberRepository FamilyMemberRepository, uuidGenerator UUIDGenerator) *UseCase {
	return &UseCase{
		taskRepository:         taskRepository,
		familyMemberRepository: familyMemberRepository,
		uuidGenerator:          uuidGenerator,
	}
}

func (u *UseCase) Add(ctx context.Context, model entity.TaskAddRequest) (entity.Task, error) {
	parentMember, err := u.familyMemberRepository.GetByUser(ctx, model.UserID)
	if err != nil {
		return entity.Task{}, fmt.Errorf("can't get family member: %w", err)
	}

	if parentMember.Role != entity.FamilyRoleParent {
		return entity.Task{}, domain.ForbiddenError()
	}

	createModel := entity.TaskCreate{
		ID:          entity.TaskID(u.uuidGenerator.GenerateV7()),
		Title:       model.Title,
		Description: model.Description,
		Points:      model.Points,
		FamilyID:    parentMember.FamilyID,
	}

	err = u.taskRepository.Create(ctx, createModel)
	if err != nil {
		return entity.Task{}, fmt.Errorf("can't create task: %w", err)
	}

	return entity.Task{
		ID:          createModel.ID,
		Title:       createModel.Title,
		Description: createModel.Description,
		Points:      createModel.Points,
		FamilyID:    createModel.FamilyID,
	}, nil
}
