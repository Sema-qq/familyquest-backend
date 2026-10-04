package list

import (
	"context"
	"fmt"

	"familyquest-backend/internal/domain/entity"
)

type SeasonRepository interface {
	ListByFamily(ctx context.Context, familyID entity.FamilyID) ([]entity.Season, error)
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

func (u *UseCase) List(ctx context.Context, userID entity.UserID) ([]entity.Season, error) {
	member, err := u.familyMemberRepository.GetByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("can't get family member: %w", err)
	}

	seasons, err := u.seasonRepository.ListByFamily(ctx, member.FamilyID)
	if err != nil {
		return nil, fmt.Errorf("can't list seasons: %w", err)
	}

	return seasons, nil
}
