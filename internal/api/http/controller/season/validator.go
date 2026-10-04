package season

import (
	"errors"
	"strings"
	"time"
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
	if endsAt.Before(startsAt) {
		return errors.New(`the "ends_at" field must be greater than or equal to "starts_at"`)
	}

	return nil
}
