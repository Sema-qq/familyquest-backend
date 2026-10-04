package get

import (
	"context"
	"fmt"

	"familyquest-backend/internal/domain"
	"familyquest-backend/internal/domain/entity"
)

type SeasonRepository interface {
	GetByID(ctx context.Context, id entity.SeasonID) (entity.Season, error)
}

type FamilyMemberRepository interface {
	GetByUser(ctx context.Context, userID entity.UserID) (entity.FamilyMember, error)
}

type UseCase struct {
	seasonRepository       SeasonRepository
	familyMemberRepository FamilyMemberRepository
}

func New(seasonRepository SeasonRepository, familyMemberRepository FamilyMemberRepository) *UseCase {
	return &UseCase{
		seasonRepository:       seasonRepository,
		familyMemberRepository: familyMemberRepository,
	}
}

func (u *UseCase) Get(ctx context.Context, userID entity.UserID, seasonID entity.SeasonID) (entity.Season, error) {
	member, err := u.familyMemberRepository.GetByUser(ctx, userID)
	if err != nil {
		return entity.Season{}, fmt.Errorf("can't get family member: %w", err)
	}

	season, err := u.seasonRepository.GetByID(ctx, seasonID)
	if err != nil {
		return entity.Season{}, fmt.Errorf("can't get season: %w", err)
	}

	if season.FamilyID != member.FamilyID {
		return entity.Season{}, domain.NotFound("season not found")
	}

	return season, nil
}
