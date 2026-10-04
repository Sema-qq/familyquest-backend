package task

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

func (r *Repository) Create(ctx context.Context, model entity.TaskCreate) error {
	_, err := r.db.Exec(
		ctx,
		createSQL,
		model.ID.UUID(),
		model.Title,
		model.Description,
		model.Points,
		model.FamilyID.UUID(),
	)
	if err != nil {
		return fmt.Errorf("failed create task: %w", err)
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, id entity.TaskID) (entity.Task, error) {
	task, err := db.QueryRow[task](ctx, r.db, getByIDSQL, id.UUID())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Task{}, domain.NotFound("task not found")
		}

		return entity.Task{}, fmt.Errorf("failed get task by id: %w", err)
	}

	return task.toEntity(), nil
}

func (r *Repository) ListByFamily(ctx context.Context, familyID entity.FamilyID) ([]entity.Task, error) {
	tasks, err := db.Query[task](ctx, r.db, listByFamilySQL, familyID.UUID())
	if err != nil {
		return nil, fmt.Errorf("failed list tasks by family: %w", err)
	}

	return taskList(tasks).toEntities(), nil
}
