package seasontaskslot

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"time"

	"familyquest-backend/internal/domain"
	"familyquest-backend/internal/domain/entity"
	"familyquest-backend/pkg/db"

	"github.com/jackc/pgx/v5"
)

//go:embed sqls/create.sql
var createSQL string

//go:embed sqls/count_by_season.sql
var countBySeasonSQL string

//go:embed sqls/find_available.sql
var findAvailableSQL string

//go:embed sqls/list_by_season.sql
var listBySeasonSQL string

type count struct {
	Count int64 `db:"count"`
}

type Repository struct {
	db db.Conn
}

func NewRepository(db db.Conn) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateMany(ctx context.Context, slots []entity.SeasonTaskSlotCreate) error {
	for _, slot := range slots {
		_, err := r.db.Exec(
			ctx,
			createSQL,
			slot.ID.UUID(),
			slot.SeasonTaskID.UUID(),
			slot.Position,
			slot.AvailableFrom,
			slot.AvailableUntil,
		)
		if err != nil {
			return fmt.Errorf("failed create season task slot: %w", err)
		}
	}

	return nil
}

func (r *Repository) CountBySeason(ctx context.Context, seasonID entity.SeasonID) (int64, error) {
	result, err := db.QueryRow[count](ctx, r.db, countBySeasonSQL, seasonID.UUID())
	if err != nil {
		return 0, fmt.Errorf("failed count season task slots by season: %w", err)
	}

	return result.Count, nil
}

func (r *Repository) FindAvailable(
	ctx context.Context,
	seasonTaskID entity.SeasonTaskID,
	memberID entity.FamilyMemberID,
	performedOn time.Time,
	shared bool,
) (entity.SeasonTaskSlot, error) {
	slot, err := db.QueryRow[seasonTaskSlot](
		ctx,
		r.db,
		findAvailableSQL,
		seasonTaskID.UUID(),
		memberID.UUID(),
		performedOn,
		shared,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.SeasonTaskSlot{}, domain.Conflict("no available task slot")
		}

		return entity.SeasonTaskSlot{}, fmt.Errorf("failed find available season task slot: %w", err)
	}

	return slot.toEntity(), nil
}

func (r *Repository) ListBySeason(ctx context.Context, seasonID entity.SeasonID) ([]entity.SeasonTaskSlot, error) {
	slots, err := db.Query[seasonTaskSlot](ctx, r.db, listBySeasonSQL, seasonID.UUID())
	if err != nil {
		return nil, fmt.Errorf("failed list season task slots by season: %w", err)
	}

	return seasonTaskSlotList(slots).toEntities(), nil
}
