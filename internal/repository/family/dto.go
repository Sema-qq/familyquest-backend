package family

import (
	"familyquest-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type family struct {
	ID       uuid.UUID `db:"id"`
	Name     string    `db:"name"`
	Timezone string    `db:"timezone"`
	OwnerID  uuid.UUID `db:"owner_id"`
}

func (f family) toEntity() entity.Family {
	return entity.Family{
		ID:       entity.FamilyID(f.ID),
		Name:     f.Name,
		Timezone: f.Timezone,
		OwnerID:  entity.UserID(f.OwnerID),
	}
}

type familyWithRole struct {
	ID       uuid.UUID `db:"id"`
	Name     string    `db:"name"`
	Timezone string    `db:"timezone"`
	OwnerID  uuid.UUID `db:"owner_id"`
	Role     string    `db:"role"`
}

func (f familyWithRole) toEntity() entity.FamilyWithRole {
	return entity.FamilyWithRole{
		Family: entity.Family{
			ID:       entity.FamilyID(f.ID),
			Name:     f.Name,
			Timezone: f.Timezone,
			OwnerID:  entity.UserID(f.OwnerID),
		},
		Role: entity.FamilyRole(f.Role),
	}
}
