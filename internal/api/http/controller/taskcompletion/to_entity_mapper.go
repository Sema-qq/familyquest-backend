package taskcompletion

import (
	"strings"
	"time"

	"familyquest-backend/internal/domain/entity"
)

const dateLayout = "2006-01-02"

type toEntityMapper struct{}

func newToEntityMapper() toEntityMapper {
	return toEntityMapper{}
}

func (m toEntityMapper) mapSubmitRequest(
	req submitRequest,
	userID entity.UserID,
	seasonTaskID entity.SeasonTaskID,
) entity.TaskCompletionSubmitRequest {
	performedOn, _ := time.Parse(dateLayout, req.PerformedOn)

	return entity.TaskCompletionSubmitRequest{
		UserID:       userID,
		SeasonTaskID: seasonTaskID,
		PerformedOn:  performedOn,
		Comment:      strings.TrimSpace(req.Comment),
	}
}
