package task

import (
	"errors"
	"strings"
)

type RequestValidator struct{}

func NewRequestValidator() *RequestValidator {
	return &RequestValidator{}
}

func (v *RequestValidator) ValidateAdd(req addRequest) error {
	if strings.TrimSpace(req.Title) == "" {
		return errors.New(`the "title" field is required`)
	}

	if req.Points <= 0 {
		return errors.New(`the "points" field must be greater than 0`)
	}

	return nil
}
