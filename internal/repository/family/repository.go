package family

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"familyquest-backend/internal/domain"
	"familyquest-backend/internal/domain/entity"
	"familyquest-backend/pkg/db"

	"github.com/jackc/pgx/v5"
)

//go:embed sqls/create.sql
var createSQL string

//go:embed sqls/get_by_user.sql
var getByUserSQL string

type Repository struct {
	db db.Conn
}

func NewRepository(db db.Conn) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, family entity.FamilyCreate) error {
	_, err := r.db.Exec(
		ctx,
		createSQL,
		family.ID.UUID(),
		family.Name,
		family.Timezone,
		family.OwnerID.UUID(),
	)
	if err != nil {
		return fmt.Errorf("failed create family: %w", err)
	}

	return nil
}

func (r *Repository) GetByUser(ctx context.Context, userID entity.UserID) (entity.FamilyWithRole, error) {
	family, err := db.QueryRow[familyWithRole](ctx, r.db, getByUserSQL, userID.UUID())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.FamilyWithRole{}, domain.NotFound("family not found")
		}

		return entity.FamilyWithRole{}, fmt.Errorf("failed get family by user: %w", err)
	}

	return family.toEntity(), nil
}
