package season

import (
	"encoding/json"
	"strings"
	"time"

	"familyquest-backend/internal/domain/entity"

	"github.com/google/uuid"
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

func (m toEntityMapper) mapAddTaskRequest(
	req addTaskRequest,
	userID entity.UserID,
	seasonID entity.SeasonID,
) entity.SeasonTaskAddRequest {
	taskID, _ := uuid.Parse(req.TaskID)

	return entity.SeasonTaskAddRequest{
		UserID:            userID,
		SeasonID:          seasonID,
		TaskID:            entity.TaskID(taskID),
		ParticipationMode: entity.TaskParticipationMode(req.ParticipationMode),
		Schedule:          m.mapSchedule(req.Schedule),
	}
}

func (m toEntityMapper) mapSchedule(raw json.RawMessage) entity.SeasonTaskSchedule {
	var scheduleType scheduleTypeRequest
	_ = json.Unmarshal(raw, &scheduleType)

	switch entity.SeasonTaskScheduleType(scheduleType.Type) {
	case entity.SeasonTaskScheduleTypeDates:
		var req datesScheduleRequest
		_ = json.Unmarshal(raw, &req)

		dates := make([]time.Time, len(req.Dates))
		for i, date := range req.Dates {
			dates[i], _ = time.Parse(dateLayout, date)
		}

		return entity.SeasonTaskSchedule{
			Type:  entity.SeasonTaskScheduleTypeDates,
			Dates: dates,
		}
	case entity.SeasonTaskScheduleTypeQuota:
		var req quotaScheduleRequest
		_ = json.Unmarshal(raw, &req)

		return entity.SeasonTaskSchedule{
			Type:  entity.SeasonTaskScheduleTypeQuota,
			Count: req.Count,
		}
	default:
		return entity.SeasonTaskSchedule{}
	}
}
