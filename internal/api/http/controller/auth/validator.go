package auth

import (
	"errors"
	"strings"
)

type RequestValidator struct{}

func NewRequestValidator() *RequestValidator {
	return &RequestValidator{}
}

func (v *RequestValidator) validateRegister(req registerRequest) error {
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

func (v *RequestValidator) validateLogin(req loginRequest) error {
	if strings.TrimSpace(req.Login) == "" {
		return errors.New(`the "login" field is required`)
	}
	if req.Password == "" {
		return errors.New(`the "password" field is required`)
	}

	return nil
}
