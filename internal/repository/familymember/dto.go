package familymember

import (
	"familyquest-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type familyMember struct {
	ID       uuid.UUID `db:"id"`
	FamilyID uuid.UUID `db:"family_id"`
	UserID   uuid.UUID `db:"user_id"`
	Role     string    `db:"role"`
}

func (m familyMember) toEntity() entity.FamilyMember {
	return entity.FamilyMember{
		ID:       entity.FamilyMemberID(m.ID),
		FamilyID: entity.FamilyID(m.FamilyID),
		UserID:   entity.UserID(m.UserID),
		Role:     entity.FamilyRole(m.Role),
	}
}

type familyMemberProfile struct {
	ID          uuid.UUID `db:"id"`
	FamilyID    uuid.UUID `db:"family_id"`
	UserID      uuid.UUID `db:"user_id"`
	Login       string    `db:"login"`
	DisplayName string    `db:"display_name"`
	Role        string    `db:"role"`
}

func (m familyMemberProfile) toEntity() entity.FamilyMemberProfile {
	return entity.FamilyMemberProfile{
		ID:          entity.FamilyMemberID(m.ID),
		FamilyID:    entity.FamilyID(m.FamilyID),
		UserID:      entity.UserID(m.UserID),
		Login:       m.Login,
		DisplayName: m.DisplayName,
		Role:        entity.FamilyRole(m.Role),
	}
}

type familyMemberProfileList []familyMemberProfile

func (list familyMemberProfileList) toEntities() []entity.FamilyMemberProfile {
	result := make([]entity.FamilyMemberProfile, len(list))
	for i, item := range list {
		result[i] = item.toEntity()
	}

	return result
}
