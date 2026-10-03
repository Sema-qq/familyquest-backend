package auth

import (
	"time"

	"familyquest-backend/internal/domain/entity"
)

type toProtocolMapper struct{}

func newToProtocolMapper() toProtocolMapper {
	return toProtocolMapper{}
}

func (m toProtocolMapper) mapRegisterResponse(user entity.User) registerResponse {
	return registerResponse{
		ID:          user.ID.String(),
		Login:       user.Login,
		DisplayName: user.DisplayName,
	}
}

func (m toProtocolMapper) mapLoginResponse(result entity.AuthResult) loginResponse {
	return loginResponse{
		Token:     result.AccessToken,
		ExpiresAt: result.ExpiresAt.Format(time.RFC3339),
	}
}
