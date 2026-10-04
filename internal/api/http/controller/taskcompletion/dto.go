package taskcompletion

type submitRequest struct {
	PerformedOn string `json:"performed_on"`
	Comment     string `json:"comment"`
}

type rejectRequest struct {
	Comment string `json:"comment"`
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

type completionViewResponse struct {
	ID            string  `json:"id"`
	MemberID      string  `json:"member_id"`
	DisplayName   string  `json:"display_name"`
	SlotID        string  `json:"slot_id"`
	SeasonTaskID  string  `json:"season_task_id"`
	TaskTitle     string  `json:"task_title"`
	PerformedOn   string  `json:"performed_on"`
	Status        string  `json:"status"`
	Points        int64   `json:"points"`
	ChildComment  string  `json:"child_comment"`
	ParentComment string  `json:"parent_comment"`
	ReviewedBy    *string `json:"reviewed_by"`
	ReviewedAt    *string `json:"reviewed_at"`
}

type completionsResponse struct {
	Items []completionViewResponse `json:"items"`
}
