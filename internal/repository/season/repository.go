package season

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

//go:embed sqls/get_by_id.sql
var getByIDSQL string

//go:embed sqls/list_by_family.sql
var listByFamilySQL string

type Repository struct {
	db db.Conn
}

func NewRepository(db db.Conn) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, season entity.SeasonCreate) error {
	_, err := r.db.Exec(
		ctx,
		createSQL,
		season.ID.UUID(),
		season.FamilyID.UUID(),
		season.Title,
		season.StartsAt,
		season.EndsAt,
		season.Timezone,
	)
	if err != nil {
		return fmt.Errorf("failed create season: %w", err)
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, id entity.SeasonID) (entity.Season, error) {
	season, err := db.QueryRow[season](ctx, r.db, getByIDSQL, id.UUID())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Season{}, domain.NotFound("season not found")
		}

		return entity.Season{}, fmt.Errorf("failed get season by id: %w", err)
	}

	return season.toEntity(), nil
}

func (r *Repository) ListByFamily(ctx context.Context, familyID entity.FamilyID) ([]entity.Season, error) {
	seasons, err := db.Query[season](ctx, r.db, listByFamilySQL, familyID.UUID())
	if err != nil {
		return nil, fmt.Errorf("failed list seasons by family: %w", err)
	}

	return seasonList(seasons).toEntities(), nil
}
