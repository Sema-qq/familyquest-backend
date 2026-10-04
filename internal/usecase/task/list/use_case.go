package list

import (
	"context"
	"fmt"

	"familyquest-backend/internal/domain/entity"
)

type TaskRepository interface {
	ListByFamily(ctx context.Context, familyID entity.FamilyID) ([]entity.Task, error)
}

type FamilyMemberRepository interface {
	GetByUser(ctx context.Context, userID entity.UserID) (entity.FamilyMember, error)
}

type UseCase struct {
	taskRepository         TaskRepository
	familyMemberRepository FamilyMemberRepository
}

func New(taskRepository TaskRepository, familyMemberRepository FamilyMemberRepository) *UseCase {
	return &UseCase{
		taskRepository:         taskRepository,
		familyMemberRepository: familyMemberRepository,
	}
}

func (u *UseCase) List(ctx context.Context, userID entity.UserID) ([]entity.Task, error) {
	member, err := u.familyMemberRepository.GetByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("can't get family member: %w", err)
	}

	tasks, err := u.taskRepository.ListByFamily(ctx, member.FamilyID)
	if err != nil {
		return nil, fmt.Errorf("can't list tasks: %w", err)
	}

	return tasks, nil
}
