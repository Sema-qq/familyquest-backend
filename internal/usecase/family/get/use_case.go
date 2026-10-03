package get

import (
	"context"
	"fmt"

	"familyquest-backend/internal/domain/entity"
)

type FamilyRepository interface {
	GetByUser(ctx context.Context, userID entity.UserID) (entity.FamilyWithRole, error)
}

type UseCase struct {
	familyRepository FamilyRepository
}

func New(familyRepository FamilyRepository) *UseCase {
	return &UseCase{
		familyRepository: familyRepository,
	}
}

func (u *UseCase) Get(ctx context.Context, userID entity.UserID) (entity.FamilyWithRole, error) {
	family, err := u.familyRepository.GetByUser(ctx, userID)
	if err != nil {
		return entity.FamilyWithRole{}, fmt.Errorf("can't get family: %w", err)
	}

	return family, nil
}
