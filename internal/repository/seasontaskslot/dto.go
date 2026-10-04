package seasontaskslot

import (
	"time"

	"familyquest-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type seasonTaskSlot struct {
	ID             uuid.UUID `db:"id"`
	SeasonTaskID   uuid.UUID `db:"season_task_id"`
	Position       int64     `db:"position"`
	AvailableFrom  time.Time `db:"available_from"`
	AvailableUntil time.Time `db:"available_until"`
}

func (s seasonTaskSlot) toEntity() entity.SeasonTaskSlot {
	return entity.SeasonTaskSlot{
		ID:             entity.SeasonTaskSlotID(s.ID),
		SeasonTaskID:   entity.SeasonTaskID(s.SeasonTaskID),
		Position:       s.Position,
		AvailableFrom:  s.AvailableFrom,
		AvailableUntil: s.AvailableUntil,
	}
}

type seasonTaskSlotList []seasonTaskSlot

func (list seasonTaskSlotList) toEntities() []entity.SeasonTaskSlot {
	result := make([]entity.SeasonTaskSlot, len(list))
	for i, item := range list {
		result[i] = item.toEntity()
	}

	return result
}
