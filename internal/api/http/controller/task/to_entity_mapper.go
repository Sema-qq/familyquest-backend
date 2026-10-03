package task

import (
	"familyquest-backend/internal/domain/entity"
)

type toEntityMapper struct{}

func newToEntityMapper() toEntityMapper {
	return toEntityMapper{}
}

func (m toEntityMapper) mapAddTask(request addRequest, userID entity.UserID) entity.TaskAddRequest {
	return entity.TaskAddRequest{
		UserID:      userID,
		Title:       request.Title,
		Description: request.Description,
		Points:      request.Points,
	}
}
