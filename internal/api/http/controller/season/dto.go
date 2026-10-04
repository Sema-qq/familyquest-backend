package season

import "encoding/json"

type createRequest struct {
	Title    string `json:"title"`
	StartsAt string `json:"starts_at"`
	EndsAt   string `json:"ends_at"`
}

type seasonResponse struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	StartsAt string `json:"starts_at"`
	EndsAt   string `json:"ends_at"`
	Status   string `json:"status"`
}

type seasonsResponse struct {
	Items []seasonResponse `json:"items"`
}

type seasonDetailResponse struct {
	ID       string                   `json:"id"`
	Title    string                   `json:"title"`
	StartsAt string                   `json:"starts_at"`
	EndsAt   string                   `json:"ends_at"`
	Status   string                   `json:"status"`
	Tasks    []seasonTaskFullResponse `json:"tasks"`
}

type addTaskRequest struct {
	TaskID            string          `json:"task_id"`
	ParticipationMode string          `json:"participation_mode"`
	Schedule          json.RawMessage `json:"schedule"`
}

type scheduleTypeRequest struct {
	Type string `json:"type"`
}

type datesScheduleRequest struct {
	Type  string   `json:"type"`
	Dates []string `json:"dates"`
}

type quotaScheduleRequest struct {
	Type  string `json:"type"`
	Count int64  `json:"count"`
}

type seasonTaskFullResponse struct {
	ID                string                   `json:"id"`
	TaskID            string                   `json:"task_id"`
	Title             string                   `json:"title"`
	Description       string                   `json:"description"`
	Points            int64                    `json:"points"`
	ParticipationMode string                   `json:"participation_mode"`
	Slots             []seasonTaskSlotResponse `json:"slots"`
}

type seasonTaskSlotResponse struct {
	ID             string               `json:"id"`
	Position       int64                `json:"position"`
	AvailableFrom  string               `json:"available_from"`
	AvailableUntil string               `json:"available_until"`
	Completions    []completionResponse `json:"completions"`
}

type completionResponse struct {
	ID            string  `json:"id"`
	MemberID      string  `json:"member_id"`
	SlotID        string  `json:"slot_id"`
	SeasonTaskID  string  `json:"season_task_id"`
	PerformedOn   string  `json:"performed_on"`
	Status        string  `json:"status"`
	ChildComment  string  `json:"child_comment"`
	ParentComment string  `json:"parent_comment"`
	ReviewedBy    *string `json:"reviewed_by"`
	ReviewedAt    *string `json:"reviewed_at"`
}

type resultsResponse struct {
	Items []resultResponse `json:"items"`
}

type resultResponse struct {
	MemberID      string `json:"member_id"`
	DisplayName   string `json:"display_name"`
	ApprovedCount int64  `json:"approved_count"`
	Points        int64  `json:"points"`
}
