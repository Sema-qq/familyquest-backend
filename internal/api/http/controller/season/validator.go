package season

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"familyquest-backend/internal/domain/entity"
)

type RequestValidator struct{}

func NewRequestValidator() *RequestValidator {
	return &RequestValidator{}
}

func (v *RequestValidator) ValidateCreate(req createRequest) error {
	if strings.TrimSpace(req.Title) == "" {
		return errors.New(`the "title" field is required`)
	}

	startsAt, err := time.Parse(dateLayout, req.StartsAt)
	if err != nil {
		return errors.New(`the "starts_at" field must be date in format YYYY-MM-DD`)
	}

	endsAt, err := time.Parse(dateLayout, req.EndsAt)
	if err != nil {
		return errors.New(`the "ends_at" field must be date in format YYYY-MM-DD`)
	}
	if !endsAt.After(startsAt) {
		return errors.New(`the "ends_at" field must be greater than "starts_at"`)
	}

	return nil
}

func (v *RequestValidator) ValidateAddTask(req addTaskRequest) error {
	if _, err := uuid.Parse(req.TaskID); err != nil {
		return errors.New(`the "task_id" field must be uuid`)
	}

	switch entity.TaskParticipationMode(req.ParticipationMode) {
	case entity.TaskParticipationModeEachChild, entity.TaskParticipationModeShared:
	default:
		return errors.New(`the "participation_mode" field must be one of: each_child, shared`)
	}

	if len(req.Schedule) == 0 {
		return errors.New(`the "schedule" field is required`)
	}

	var scheduleType scheduleTypeRequest
	if err := json.Unmarshal(req.Schedule, &scheduleType); err != nil {
		return fmt.Errorf(`the "schedule" field must be object: %w`, err)
	}

	switch entity.SeasonTaskScheduleType(scheduleType.Type) {
	case entity.SeasonTaskScheduleTypeDates:
		return v.validateDatesSchedule(req.Schedule)
	case entity.SeasonTaskScheduleTypeQuota:
		return v.validateQuotaSchedule(req.Schedule)
	default:
		return errors.New(`the "schedule.type" field must be one of: dates, quota`)
	}
}

func (v *RequestValidator) validateDatesSchedule(raw json.RawMessage) error {
	var req datesScheduleRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		return fmt.Errorf(`the "schedule" field must be dates schedule: %w`, err)
	}
	if len(req.Dates) == 0 {
		return errors.New(`the "schedule.dates" field is required`)
	}

	seen := make(map[string]struct{}, len(req.Dates))
	for _, date := range req.Dates {
		if _, err := time.Parse(dateLayout, date); err != nil {
			return errors.New(`the "schedule.dates" field must contain dates in format YYYY-MM-DD`)
		}
		if _, ok := seen[date]; ok {
			return errors.New(`the "schedule.dates" field must not contain duplicates`)
		}

		seen[date] = struct{}{}
	}

	return nil
}

func (v *RequestValidator) validateQuotaSchedule(raw json.RawMessage) error {
	var req quotaScheduleRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		return fmt.Errorf(`the "schedule" field must be quota schedule: %w`, err)
	}
	if req.Count < 1 || req.Count > 100 {
		return errors.New(`the "schedule.count" field must be between 1 and 100`)
	}

	return nil
}
