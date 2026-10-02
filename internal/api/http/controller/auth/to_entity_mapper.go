package auth

import (
	"strings"

	"familyquest-backend/internal/domain/entity"
)

type toEntityMapper struct{}

func newToEntityMapper() *toEntityMapper {
	return &toEntityMapper{}
}

func (m *toEntityMapper) mapUserCreateRequest(req registerRequest) entity.UserCreateRequest {
	return entity.UserCreateRequest{
		Login:       strings.TrimSpace(req.Login),
		Password:    req.Password,
		DisplayName: strings.TrimSpace(req.DisplayName),
	}
}

func (m *toEntityMapper) mapAuthRequest(req loginRequest) entity.AuthRequest {
	return entity.AuthRequest{
		Login:    strings.TrimSpace(req.Login),
		Password: req.Password,
	}
}
