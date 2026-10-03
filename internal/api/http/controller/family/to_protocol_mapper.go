package family

import "familyquest-backend/internal/domain/entity"

type toProtocolMapper struct{}

func newToProtocolMapper() toProtocolMapper {
	return toProtocolMapper{}
}

func (m toProtocolMapper) mapFamilyResponse(family entity.FamilyWithRole) familyResponse {
	return familyResponse{
		ID:       family.Family.ID.String(),
		Name:     family.Family.Name,
		Timezone: family.Family.Timezone,
		OwnerID:  family.Family.OwnerID.String(),
		Role:     string(family.Role),
	}
}

func (m toProtocolMapper) mapFamilyMemberResponse(member entity.FamilyMemberProfile) memberResponse {
	return memberResponse{
		ID:          member.ID.String(),
		UserID:      member.UserID.String(),
		Login:       member.Login,
		DisplayName: member.DisplayName,
		Role:        string(member.Role),
	}
}

func (m toProtocolMapper) mapFamilyMembersResponse(members []entity.FamilyMemberProfile) membersResponse {
	items := make([]memberResponse, len(members))
	for i, member := range members {
		items[i] = m.mapFamilyMemberResponse(member)
	}

	return membersResponse{
		Items: items,
	}
}
