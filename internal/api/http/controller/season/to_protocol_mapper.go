package season

import "familyquest-backend/internal/domain/entity"

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

func (m toProtocolMapper) mapSeasonDetailResponse(season entity.Season) seasonDetailResponse {
	return seasonDetailResponse{
		ID:       season.ID.String(),
		Title:    season.Title,
		StartsAt: season.StartsAt.Format(dateLayout),
		EndsAt:   season.EndsAt.Format(dateLayout),
		Status:   string(season.Status),
		Tasks:    make([]seasonTaskResponse, 0),
	}
}
