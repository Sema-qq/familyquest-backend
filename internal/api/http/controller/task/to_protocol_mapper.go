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
