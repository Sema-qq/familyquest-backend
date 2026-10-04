package task

import (
	"familyquest-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type task struct {
	ID          uuid.UUID `db:"id"`
	FamilyID    uuid.UUID `db:"family_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Points      int64     `db:"points"`
}

func (t task) toEntity() entity.Task {
	return entity.Task{
		ID:          entity.TaskID(t.ID),
		FamilyID:    entity.FamilyID(t.FamilyID),
		Title:       t.Title,
		Description: t.Description,
		Points:      t.Points,
	}
}

type taskList []task

func (list taskList) toEntities() []entity.Task {
	result := make([]entity.Task, len(list))
	for i, item := range list {
		result[i] = item.toEntity()
	}

	return result
}
