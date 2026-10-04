package season

import (
	"time"

	"familyquest-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type season struct {
	ID          uuid.UUID  `db:"id"`
	FamilyID    uuid.UUID  `db:"family_id"`
	Title       string     `db:"title"`
	StartsAt    time.Time  `db:"starts_at"`
	EndsAt      time.Time  `db:"ends_at"`
	Timezone    string     `db:"timezone"`
	Status      string     `db:"status"`
	CompletedAt *time.Time `db:"completed_at"`
}

func (s season) toEntity() entity.Season {
	return entity.Season{
		ID:        entity.SeasonID(s.ID),
		FamilyID:  entity.FamilyID(s.FamilyID),
		Title:     s.Title,
		StartsAt:  s.StartsAt,
		EndsAt:    s.EndsAt,
		Timezone:  s.Timezone,
		Status:    entity.SeasonStatus(s.Status),
		Completed: s.CompletedAt != nil,
	}
}

type seasonList []season

func (list seasonList) toEntities() []entity.Season {
	result := make([]entity.Season, len(list))
	for i, item := range list {
		result[i] = item.toEntity()
	}

	return result
}
