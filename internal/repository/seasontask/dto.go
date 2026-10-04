package seasontask

import (
	"time"

	"familyquest-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type seasonTaskWithSeason struct {
	ID                uuid.UUID `db:"id"`
	SeasonID          uuid.UUID `db:"season_id"`
	TaskID            uuid.UUID `db:"task_id"`
	Title             string    `db:"title"`
	Description       string    `db:"description"`
	Points            int64     `db:"points"`
	ParticipationMode string    `db:"participation_mode"`
	SeasonFamilyID    uuid.UUID `db:"season_family_id"`
	SeasonStartsAt    time.Time `db:"season_starts_at"`
	SeasonEndsAt      time.Time `db:"season_ends_at"`
	SeasonTimezone    string    `db:"season_timezone"`
	SeasonStatus      string    `db:"season_status"`
}

type seasonTask struct {
	ID                uuid.UUID `db:"id"`
	SeasonID          uuid.UUID `db:"season_id"`
	TaskID            uuid.UUID `db:"task_id"`
	Title             string    `db:"title"`
	Description       string    `db:"description"`
	Points            int64     `db:"points"`
	ParticipationMode string    `db:"participation_mode"`
}

func (s seasonTask) toEntity() entity.SeasonTask {
	return entity.SeasonTask{
		ID:                entity.SeasonTaskID(s.ID),
		SeasonID:          entity.SeasonID(s.SeasonID),
		TaskID:            entity.TaskID(s.TaskID),
		Title:             s.Title,
		Description:       s.Description,
		Points:            s.Points,
		ParticipationMode: entity.TaskParticipationMode(s.ParticipationMode),
	}
}

type seasonTaskList []seasonTask

func (list seasonTaskList) toEntities() []entity.SeasonTask {
	result := make([]entity.SeasonTask, len(list))
	for i, item := range list {
		result[i] = item.toEntity()
	}

	return result
}

func (s seasonTaskWithSeason) toEntity() entity.SeasonTaskWithSeason {
	return entity.SeasonTaskWithSeason{
		SeasonTask: entity.SeasonTask{
			ID:                entity.SeasonTaskID(s.ID),
			SeasonID:          entity.SeasonID(s.SeasonID),
			TaskID:            entity.TaskID(s.TaskID),
			Title:             s.Title,
			Description:       s.Description,
			Points:            s.Points,
			ParticipationMode: entity.TaskParticipationMode(s.ParticipationMode),
		},
		SeasonFamilyID: entity.FamilyID(s.SeasonFamilyID),
		SeasonStartsAt: s.SeasonStartsAt,
		SeasonEndsAt:   s.SeasonEndsAt,
		SeasonTimezone: s.SeasonTimezone,
		SeasonStatus:   entity.SeasonStatus(s.SeasonStatus),
	}
}
