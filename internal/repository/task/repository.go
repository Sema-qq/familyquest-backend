package task

import (
	"context"
	_ "embed"
	"fmt"

	"familyquest-backend/internal/domain/entity"
	"familyquest-backend/pkg/db"
)

//go:embed sqls/create.sql
var createSQL string

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
