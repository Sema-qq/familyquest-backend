package taskcompletion

import (
	"time"

	"familyquest-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type toProtocolMapper struct{}

func newToProtocolMapper() toProtocolMapper {
	return toProtocolMapper{}
}

func (m toProtocolMapper) mapCompletionResponse(completion entity.TaskCompletion) completionResponse {
	var reviewedBy *string
	if completion.ReviewedBy.UUID() != uuid.Nil {
		value := completion.ReviewedBy.String()
		reviewedBy = &value
	}

	var reviewedAt *string
	if completion.ReviewedAt != nil {
		value := completion.ReviewedAt.Format(time.RFC3339)
		reviewedAt = &value
	}

	return completionResponse{
		ID:            completion.ID.String(),
		MemberID:      completion.MemberID.String(),
		SlotID:        completion.SlotID.String(),
		SeasonTaskID:  completion.SeasonTaskID.String(),
		PerformedOn:   completion.PerformedOn.Format(dateLayout),
		Status:        string(completion.Status),
		ChildComment:  completion.ChildComment,
		ParentComment: completion.ParentComment,
		ReviewedBy:    reviewedBy,
		ReviewedAt:    reviewedAt,
	}
}

func (m toProtocolMapper) mapCompletionsResponse(completions []entity.TaskCompletionView) completionsResponse {
	items := make([]completionViewResponse, len(completions))
	for i, completion := range completions {
		items[i] = m.mapCompletionViewResponse(completion)
	}

	return completionsResponse{
		Items: items,
	}
}

func (m toProtocolMapper) mapCompletionViewResponse(completion entity.TaskCompletionView) completionViewResponse {
	var reviewedBy *string
	if completion.ReviewedBy.UUID() != uuid.Nil {
		value := completion.ReviewedBy.String()
		reviewedBy = &value
	}

	var reviewedAt *string
	if completion.ReviewedAt != nil {
		value := completion.ReviewedAt.Format(time.RFC3339)
		reviewedAt = &value
	}

	return completionViewResponse{
		ID:            completion.ID.String(),
		MemberID:      completion.MemberID.String(),
		DisplayName:   completion.DisplayName,
		SlotID:        completion.SlotID.String(),
		SeasonTaskID:  completion.SeasonTaskID.String(),
		TaskTitle:     completion.TaskTitle,
		PerformedOn:   completion.PerformedOn.Format(dateLayout),
		Status:        string(completion.Status),
		Points:        completion.Points,
		ChildComment:  completion.ChildComment,
		ParentComment: completion.ParentComment,
		ReviewedBy:    reviewedBy,
		ReviewedAt:    reviewedAt,
	}
}
