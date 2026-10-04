package seasontask

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"familyquest-backend/internal/domain"
	"familyquest-backend/internal/domain/entity"
	"familyquest-backend/pkg/db"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

//go:embed sqls/create.sql
var createSQL string

//go:embed sqls/count_by_season.sql
var countBySeasonSQL string

//go:embed sqls/get_by_id.sql
var getByIDSQL string

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

func (r *Repository) Create(ctx context.Context, task entity.SeasonTaskCreate) error {
	_, err := r.db.Exec(
		ctx,
		createSQL,
		task.ID.UUID(),
		task.SeasonID.UUID(),
		task.TaskID.UUID(),
		task.Title,
		task.Description,
		task.Points,
		string(task.ParticipationMode),
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.Conflict("task already added to season")
		}

		return fmt.Errorf("failed create season task: %w", err)
	}

	return nil
}

func (r *Repository) CountBySeason(ctx context.Context, seasonID entity.SeasonID) (int64, error) {
	result, err := db.QueryRow[count](ctx, r.db, countBySeasonSQL, seasonID.UUID())
	if err != nil {
		return 0, fmt.Errorf("failed count season tasks by season: %w", err)
	}

	return result.Count, nil
}

func (r *Repository) GetByID(ctx context.Context, id entity.SeasonTaskID) (entity.SeasonTaskWithSeason, error) {
	seasonTask, err := db.QueryRow[seasonTaskWithSeason](ctx, r.db, getByIDSQL, id.UUID())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.SeasonTaskWithSeason{}, domain.NotFound("season task not found")
		}

		return entity.SeasonTaskWithSeason{}, fmt.Errorf("failed get season task by id: %w", err)
	}

	return seasonTask.toEntity(), nil
}

func (r *Repository) ListBySeason(ctx context.Context, seasonID entity.SeasonID) ([]entity.SeasonTask, error) {
	tasks, err := db.Query[seasonTask](ctx, r.db, listBySeasonSQL, seasonID.UUID())
	if err != nil {
		return nil, fmt.Errorf("failed list season tasks by season: %w", err)
	}

	return seasonTaskList(tasks).toEntities(), nil
}
