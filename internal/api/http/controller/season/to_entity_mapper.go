package season

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

func (m toEntityMapper) mapCreateRequest(req createRequest, userID entity.UserID) entity.SeasonCreateRequest {
	startsAt, _ := time.Parse(dateLayout, req.StartsAt)
	endsAt, _ := time.Parse(dateLayout, req.EndsAt)

	return entity.SeasonCreateRequest{
		UserID:   userID,
		Title:    strings.TrimSpace(req.Title),
		StartsAt: startsAt,
		EndsAt:   endsAt,
	}
}
