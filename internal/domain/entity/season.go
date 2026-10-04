package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type SeasonID uuid.UUID

func (id SeasonID) String() string {
	return id.UUID().String()
}

func (id SeasonID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

type SeasonStatus string

const (
	SeasonStatusDraft     SeasonStatus = "draft"
	SeasonStatusActive    SeasonStatus = "active"
	SeasonStatusCompleted SeasonStatus = "completed"
)

type Season struct {
	ID        SeasonID
	FamilyID  FamilyID
	Title     string
	StartsAt  time.Time
	EndsAt    time.Time
	Timezone  string
	Status    SeasonStatus
	Completed bool
}

type SeasonDetail struct {
	Season Season
	Tasks  []SeasonTaskDetail
}

type SeasonCreate struct {
	ID       SeasonID
	FamilyID FamilyID
	Title    string
	StartsAt time.Time
	EndsAt   time.Time
	Timezone string
}

type SeasonCreateRequest struct {
	UserID   UserID
	Title    string
	StartsAt time.Time
	EndsAt   time.Time
}

type SeasonUpdate struct {
	ID          SeasonID
	Status      null.Value[SeasonStatus]
	CompletedAt null.Value[time.Time]
}

type SeasonTaskID uuid.UUID

func (id SeasonTaskID) String() string {
	return id.UUID().String()
}

func (id SeasonTaskID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

type SeasonTask struct {
	ID                SeasonTaskID
	SeasonID          SeasonID
	TaskID            TaskID
	Title             string
	Description       string
	Points            int64
	ParticipationMode TaskParticipationMode
	Slots             []SeasonTaskSlot
}

type SeasonTaskDetail struct {
	SeasonTask
	Slots []SeasonTaskSlotDetail
}

type SeasonTaskWithSeason struct {
	SeasonTask
	SeasonFamilyID FamilyID
	SeasonStartsAt time.Time
	SeasonEndsAt   time.Time
	SeasonTimezone string
	SeasonStatus   SeasonStatus
}

type SeasonTaskCreate struct {
	ID                SeasonTaskID
	SeasonID          SeasonID
	TaskID            TaskID
	Title             string
	Description       string
	Points            int64
	ParticipationMode TaskParticipationMode
}

type SeasonTaskSlotID uuid.UUID

func (id SeasonTaskSlotID) String() string {
	return id.UUID().String()
}

func (id SeasonTaskSlotID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

type SeasonTaskSlot struct {
	ID             SeasonTaskSlotID
	SeasonTaskID   SeasonTaskID
	Position       int64
	AvailableFrom  time.Time
	AvailableUntil time.Time
}

type SeasonTaskSlotDetail struct {
	SeasonTaskSlot
	Completions []TaskCompletion
}

type SeasonTaskSlotCreate struct {
	ID             SeasonTaskSlotID
	SeasonTaskID   SeasonTaskID
	Position       int64
	AvailableFrom  time.Time
	AvailableUntil time.Time
}

type SeasonTaskScheduleType string

const (
	SeasonTaskScheduleTypeDates SeasonTaskScheduleType = "dates"
	SeasonTaskScheduleTypeQuota SeasonTaskScheduleType = "quota"
)

type SeasonTaskAddRequest struct {
	UserID            UserID
	SeasonID          SeasonID
	TaskID            TaskID
	ParticipationMode TaskParticipationMode
	Schedule          SeasonTaskSchedule
}

type SeasonTaskSchedule struct {
	Type  SeasonTaskScheduleType
	Dates []time.Time
	Count int64
}
