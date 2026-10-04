package taskcompletion

import (
	"errors"
	"strings"
	"time"

	"familyquest-backend/internal/domain/entity"
)

type RequestValidator struct{}

func NewRequestValidator() *RequestValidator {
	return &RequestValidator{}
}

func (v *RequestValidator) ValidateSubmit(req submitRequest) error {
	if _, err := time.Parse(dateLayout, req.PerformedOn); err != nil {
		return errors.New(`the "performed_on" field must be date in format YYYY-MM-DD`)
	}

	return nil
}

func (v *RequestValidator) ValidateStatus(status string) (*entity.TaskCompletionStatus, error) {
	if status == "" {
		return nil, nil
	}

	result := entity.TaskCompletionStatus(status)
	switch result {
	case entity.TaskCompletionStatusSubmitted, entity.TaskCompletionStatusApproved, entity.TaskCompletionStatusRejected:
		return &result, nil
	default:
		return nil, errors.New(`the "status" query parameter must be one of: submitted, approved, rejected`)
	}
}

func (v *RequestValidator) ValidateReject(req rejectRequest) error {
	if strings.TrimSpace(req.Comment) == "" {
		return errors.New(`the "comment" field is required`)
	}

	return nil
}
