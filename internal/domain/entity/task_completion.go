package entity

import (
	"time"

	"github.com/google/uuid"
)

type TaskCompletionID uuid.UUID

func (id TaskCompletionID) String() string {
	return id.UUID().String()
}

func (id TaskCompletionID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

type TaskCompletionStatus string

const (
	TaskCompletionStatusSubmitted TaskCompletionStatus = "submitted"
	TaskCompletionStatusApproved  TaskCompletionStatus = "approved"
	TaskCompletionStatusRejected  TaskCompletionStatus = "rejected"
)

type TaskCompletion struct {
	ID            TaskCompletionID
	MemberID      FamilyMemberID
	SlotID        SeasonTaskSlotID
	SeasonTaskID  SeasonTaskID
	PerformedOn   time.Time
	Status        TaskCompletionStatus
	ChildComment  string
	ParentComment string
	ReviewedBy    FamilyMemberID
	ReviewedAt    *time.Time
}

type TaskCompletionView struct {
	ID            TaskCompletionID
	MemberID      FamilyMemberID
	DisplayName   string
	SlotID        SeasonTaskSlotID
	SeasonTaskID  SeasonTaskID
	TaskTitle     string
	PerformedOn   time.Time
	Status        TaskCompletionStatus
	Points        int64
	ChildComment  string
	ParentComment string
	ReviewedBy    FamilyMemberID
	ReviewedAt    *time.Time
}

type TaskCompletionReview struct {
	TaskCompletion
	SeasonID          SeasonID
	SeasonFamilyID    FamilyID
	SeasonStatus      SeasonStatus
	ParticipationMode TaskParticipationMode
}

type TaskCompletionUpdateReview struct {
	ID            TaskCompletionID
	Status        TaskCompletionStatus
	ParentComment string
	ReviewedBy    FamilyMemberID
	ReviewedAt    time.Time
}

type TaskCompletionListRequest struct {
	UserID   UserID
	SeasonID SeasonID
	Status   *TaskCompletionStatus
}

type TaskCompletionCreate struct {
	ID           TaskCompletionID
	MemberID     FamilyMemberID
	SlotID       SeasonTaskSlotID
	PerformedOn  time.Time
	ChildComment string
}

type TaskCompletionResubmit struct {
	ID           TaskCompletionID
	PerformedOn  time.Time
	ChildComment string
}

type TaskCompletionSubmitRequest struct {
	UserID       UserID
	SeasonTaskID SeasonTaskID
	PerformedOn  time.Time
	Comment      string
}

type TaskCompletionReviewRequest struct {
	UserID       UserID
	CompletionID TaskCompletionID
	Comment      string
}

type SeasonResult struct {
	MemberID      FamilyMemberID
	DisplayName   string
	ApprovedCount int64
	Points        int64
}
