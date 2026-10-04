package entity

import (
	"time"

	"github.com/google/uuid"
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
