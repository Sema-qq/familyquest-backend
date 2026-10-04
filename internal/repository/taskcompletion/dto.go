package taskcompletion

import (
	"time"

	"familyquest-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type taskCompletion struct {
	ID            uuid.UUID  `db:"id"`
	MemberID      uuid.UUID  `db:"member_id"`
	SlotID        uuid.UUID  `db:"slot_id"`
	SeasonTaskID  uuid.UUID  `db:"season_task_id"`
	PerformedOn   time.Time  `db:"performed_on"`
	Status        string     `db:"status"`
	ChildComment  string     `db:"child_comment"`
	ParentComment string     `db:"parent_comment"`
	ReviewedBy    *uuid.UUID `db:"reviewed_by"`
	ReviewedAt    *time.Time `db:"reviewed_at"`
}

func (t taskCompletion) toEntity() entity.TaskCompletion {
	result := entity.TaskCompletion{
		ID:            entity.TaskCompletionID(t.ID),
		MemberID:      entity.FamilyMemberID(t.MemberID),
		SlotID:        entity.SeasonTaskSlotID(t.SlotID),
		SeasonTaskID:  entity.SeasonTaskID(t.SeasonTaskID),
		PerformedOn:   t.PerformedOn,
		Status:        entity.TaskCompletionStatus(t.Status),
		ChildComment:  t.ChildComment,
		ParentComment: t.ParentComment,
		ReviewedAt:    t.ReviewedAt,
	}
	if t.ReviewedBy != nil {
		result.ReviewedBy = entity.FamilyMemberID(*t.ReviewedBy)
	}

	return result
}

type taskCompletionView struct {
	ID            uuid.UUID  `db:"id"`
	MemberID      uuid.UUID  `db:"member_id"`
	DisplayName   string     `db:"display_name"`
	SlotID        uuid.UUID  `db:"slot_id"`
	SeasonTaskID  uuid.UUID  `db:"season_task_id"`
	TaskTitle     string     `db:"task_title"`
	PerformedOn   time.Time  `db:"performed_on"`
	Status        string     `db:"status"`
	Points        int64      `db:"points"`
	ChildComment  string     `db:"child_comment"`
	ParentComment string     `db:"parent_comment"`
	ReviewedBy    *uuid.UUID `db:"reviewed_by"`
	ReviewedAt    *time.Time `db:"reviewed_at"`
}

func (t taskCompletionView) toEntity() entity.TaskCompletionView {
	result := entity.TaskCompletionView{
		ID:            entity.TaskCompletionID(t.ID),
		MemberID:      entity.FamilyMemberID(t.MemberID),
		DisplayName:   t.DisplayName,
		SlotID:        entity.SeasonTaskSlotID(t.SlotID),
		SeasonTaskID:  entity.SeasonTaskID(t.SeasonTaskID),
		TaskTitle:     t.TaskTitle,
		PerformedOn:   t.PerformedOn,
		Status:        entity.TaskCompletionStatus(t.Status),
		Points:        t.Points,
		ChildComment:  t.ChildComment,
		ParentComment: t.ParentComment,
		ReviewedAt:    t.ReviewedAt,
	}
	if t.ReviewedBy != nil {
		result.ReviewedBy = entity.FamilyMemberID(*t.ReviewedBy)
	}

	return result
}

type taskCompletionViewList []taskCompletionView

func (list taskCompletionViewList) toEntities() []entity.TaskCompletionView {
	result := make([]entity.TaskCompletionView, len(list))
	for i, item := range list {
		result[i] = item.toEntity()
	}

	return result
}

type taskCompletionReview struct {
	ID                uuid.UUID  `db:"id"`
	MemberID          uuid.UUID  `db:"member_id"`
	SlotID            uuid.UUID  `db:"slot_id"`
	SeasonTaskID      uuid.UUID  `db:"season_task_id"`
	PerformedOn       time.Time  `db:"performed_on"`
	Status            string     `db:"status"`
	ChildComment      string     `db:"child_comment"`
	ParentComment     string     `db:"parent_comment"`
	ReviewedBy        *uuid.UUID `db:"reviewed_by"`
	ReviewedAt        *time.Time `db:"reviewed_at"`
	SeasonID          uuid.UUID  `db:"season_id"`
	SeasonFamilyID    uuid.UUID  `db:"season_family_id"`
	SeasonStatus      string     `db:"season_status"`
	ParticipationMode string     `db:"participation_mode"`
}

func (t taskCompletionReview) toEntity() entity.TaskCompletionReview {
	result := entity.TaskCompletionReview{
		TaskCompletion: entity.TaskCompletion{
			ID:            entity.TaskCompletionID(t.ID),
			MemberID:      entity.FamilyMemberID(t.MemberID),
			SlotID:        entity.SeasonTaskSlotID(t.SlotID),
			SeasonTaskID:  entity.SeasonTaskID(t.SeasonTaskID),
			PerformedOn:   t.PerformedOn,
			Status:        entity.TaskCompletionStatus(t.Status),
			ChildComment:  t.ChildComment,
			ParentComment: t.ParentComment,
			ReviewedAt:    t.ReviewedAt,
		},
		SeasonID:          entity.SeasonID(t.SeasonID),
		SeasonFamilyID:    entity.FamilyID(t.SeasonFamilyID),
		SeasonStatus:      entity.SeasonStatus(t.SeasonStatus),
		ParticipationMode: entity.TaskParticipationMode(t.ParticipationMode),
	}
	if t.ReviewedBy != nil {
		result.ReviewedBy = entity.FamilyMemberID(*t.ReviewedBy)
	}

	return result
}

type seasonResult struct {
	MemberID      uuid.UUID `db:"member_id"`
	DisplayName   string    `db:"display_name"`
	ApprovedCount int64     `db:"approved_count"`
	Points        int64     `db:"points"`
}

func (r seasonResult) toEntity() entity.SeasonResult {
	return entity.SeasonResult{
		MemberID:      entity.FamilyMemberID(r.MemberID),
		DisplayName:   r.DisplayName,
		ApprovedCount: r.ApprovedCount,
		Points:        r.Points,
	}
}

type seasonResultList []seasonResult

func (list seasonResultList) toEntities() []entity.SeasonResult {
	result := make([]entity.SeasonResult, len(list))
	for i, item := range list {
		result[i] = item.toEntity()
	}

	return result
}
