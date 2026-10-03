package user

import "familyquest-backend/internal/domain/entity"

type toProtocolMapper struct{}

func newToProtocolMapper() toProtocolMapper {
	return toProtocolMapper{}
}

func (m toProtocolMapper) mapUserResponse(user entity.User) userResponse {
	return userResponse{
		ID:          user.ID.String(),
		Login:       user.Login,
		DisplayName: user.DisplayName,
	}
}
