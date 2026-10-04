package season

import (
	"time"

	"familyquest-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type toProtocolMapper struct{}

func newToProtocolMapper() toProtocolMapper {
	return toProtocolMapper{}
}

func (m toProtocolMapper) mapSeasonResponse(season entity.Season) seasonResponse {
	return seasonResponse{
		ID:       season.ID.String(),
		Title:    season.Title,
		StartsAt: season.StartsAt.Format(dateLayout),
		EndsAt:   season.EndsAt.Format(dateLayout),
		Status:   string(season.Status),
	}
}

func (m toProtocolMapper) mapSeasonsResponse(seasons []entity.Season) seasonsResponse {
	items := make([]seasonResponse, len(seasons))
	for i, item := range seasons {
		items[i] = m.mapSeasonResponse(item)
	}

	return seasonsResponse{
		Items: items,
	}
}

func (m toProtocolMapper) mapSeasonDetailResponse(detail entity.SeasonDetail) seasonDetailResponse {
	tasks := make([]seasonTaskFullResponse, len(detail.Tasks))
	for i, task := range detail.Tasks {
		tasks[i] = m.mapSeasonTaskDetailResponse(task)
	}

	return seasonDetailResponse{
		ID:       detail.Season.ID.String(),
		Title:    detail.Season.Title,
		StartsAt: detail.Season.StartsAt.Format(dateLayout),
		EndsAt:   detail.Season.EndsAt.Format(dateLayout),
		Status:   string(detail.Season.Status),
		Tasks:    tasks,
	}
}

func (m toProtocolMapper) mapSeasonTaskResponse(task entity.SeasonTask) seasonTaskFullResponse {
	slots := make([]seasonTaskSlotResponse, len(task.Slots))
	for i, slot := range task.Slots {
		slots[i] = seasonTaskSlotResponse{
			ID:             slot.ID.String(),
			Position:       slot.Position,
			AvailableFrom:  slot.AvailableFrom.Format(dateLayout),
			AvailableUntil: slot.AvailableUntil.Format(dateLayout),
			Completions:    make([]completionResponse, 0),
		}
	}

	return seasonTaskFullResponse{
		ID:                task.ID.String(),
		TaskID:            task.TaskID.String(),
		Title:             task.Title,
		Description:       task.Description,
		Points:            task.Points,
		ParticipationMode: string(task.ParticipationMode),
		Slots:             slots,
	}
}

func (m toProtocolMapper) mapSeasonTaskDetailResponse(task entity.SeasonTaskDetail) seasonTaskFullResponse {
	slots := make([]seasonTaskSlotResponse, len(task.Slots))
	for i, slot := range task.Slots {
		completions := make([]completionResponse, len(slot.Completions))
		for j, completion := range slot.Completions {
			completions[j] = m.mapCompletionResponse(completion)
		}

		slots[i] = seasonTaskSlotResponse{
			ID:             slot.ID.String(),
			Position:       slot.Position,
			AvailableFrom:  slot.AvailableFrom.Format(dateLayout),
			AvailableUntil: slot.AvailableUntil.Format(dateLayout),
			Completions:    completions,
		}
	}

	return seasonTaskFullResponse{
		ID:                task.ID.String(),
		TaskID:            task.TaskID.String(),
		Title:             task.Title,
		Description:       task.Description,
		Points:            task.Points,
		ParticipationMode: string(task.ParticipationMode),
		Slots:             slots,
	}
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

func (m toProtocolMapper) mapResultsResponse(results []entity.SeasonResult) resultsResponse {
	items := make([]resultResponse, len(results))
	for i, result := range results {
		items[i] = resultResponse{
			MemberID:      result.MemberID.String(),
			DisplayName:   result.DisplayName,
			ApprovedCount: result.ApprovedCount,
			Points:        result.Points,
		}
	}

	return resultsResponse{
		Items: items,
	}
}
