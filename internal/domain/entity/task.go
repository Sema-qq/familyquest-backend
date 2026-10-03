package entity

import (
	"github.com/google/uuid"
)

type TaskID uuid.UUID

func (id TaskID) String() string {
	return uuid.UUID(id).String()
}

func (id TaskID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

type Task struct {
	ID          TaskID
	Title       string
	Description string
	FamilyID    FamilyID
	Points      int64
}

type TaskCreate struct {
	ID          TaskID
	Title       string
	Description string
	FamilyID    FamilyID
	Points      int64
}

type TaskAddRequest struct {
	UserID      UserID
	Title       string
	Description string
	Points      int64
}
