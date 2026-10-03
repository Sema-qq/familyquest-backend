package family

import (
	"strings"

	"familyquest-backend/internal/domain/entity"
)

type toEntityMapper struct{}

func newToEntityMapper() toEntityMapper {
	return toEntityMapper{}
}

func (m toEntityMapper) mapFamilyCreateRequest(req createRequest, ownerID entity.UserID) entity.FamilyCreateRequest {
	return entity.FamilyCreateRequest{
		Name:     strings.TrimSpace(req.Name),
		Timezone: strings.TrimSpace(req.Timezone),
		OwnerID:  ownerID,
	}
}

func (m toEntityMapper) mapFamilyMemberCreateRequest(
	req createMemberRequest,
	parentID entity.UserID,
) entity.FamilyMemberCreateRequest {
	return entity.FamilyMemberCreateRequest{
		Login:       strings.TrimSpace(req.Login),
		Password:    req.Password,
		DisplayName: strings.TrimSpace(req.DisplayName),
		ParentID:    parentID,
	}
}
