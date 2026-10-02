package health

import (
	"context"

	"familyquest-backend/pkg/db"
)

type Repository struct {
	db db.Conn
}

func NewRepository(db db.Conn) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Ping(ctx context.Context) error {
	return r.db.Ping(ctx)
}
