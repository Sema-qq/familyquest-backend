package family

import (
	"errors"
	"strings"
)

type RequestValidator struct{}

func NewRequestValidator() *RequestValidator {
	return &RequestValidator{}
}

func (v *RequestValidator) ValidateCreate(req createRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return errors.New(`the "name" field is required`)
	}

	if strings.TrimSpace(req.Timezone) == "" {
		return errors.New(`the "timezone" field is required`)
	}

	return nil
}

func (v *RequestValidator) ValidateCreateMember(req createMemberRequest) error {
	if strings.TrimSpace(req.Login) == "" {
		return errors.New(`the "login" field is required`)
	}

	if req.Password == "" {
		return errors.New(`the "password" field is required`)
	}

	if strings.TrimSpace(req.DisplayName) == "" {
		return errors.New(`the "display_name" field is required`)
	}

	return nil
}
