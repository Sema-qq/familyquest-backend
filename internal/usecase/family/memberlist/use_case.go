package memberlist

import (
	"context"
	"fmt"

	"familyquest-backend/internal/domain/entity"
)

type FamilyMemberRepository interface {
	GetByUser(ctx context.Context, userID entity.UserID) (entity.FamilyMember, error)
	ListByFamily(ctx context.Context, familyID entity.FamilyID) ([]entity.FamilyMemberProfile, error)
}

type UseCase struct {
	familyMemberRepository FamilyMemberRepository
}

func New(familyMemberRepository FamilyMemberRepository) *UseCase {
	return &UseCase{
		familyMemberRepository: familyMemberRepository,
	}
}

func (u *UseCase) List(ctx context.Context, userID entity.UserID) ([]entity.FamilyMemberProfile, error) {
	member, err := u.familyMemberRepository.GetByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("can't get family member: %w", err)
	}

	members, err := u.familyMemberRepository.ListByFamily(ctx, member.FamilyID)
	if err != nil {
		return nil, fmt.Errorf("can't list family members: %w", err)
	}

	return members, nil
}
