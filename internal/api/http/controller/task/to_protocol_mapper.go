package task

import "familyquest-backend/internal/domain/entity"

type toProtocolMapper struct{}

func newToProtocolMapper() toProtocolMapper {
	return toProtocolMapper{}
}

func (m toProtocolMapper) mapAddResponse(item entity.Task) addResponse {
	return addResponse{
		ID:          item.ID.String(),
		Title:       item.Title,
		Description: item.Description,
		Points:      item.Points,
	}
}

func (m toProtocolMapper) mapTasksResponse(tasks []entity.Task) tasksResponse {
	items := make([]addResponse, len(tasks))
	for i, item := range tasks {
		items[i] = m.mapAddResponse(item)
	}

	return tasksResponse{
		Items: items,
	}
}
