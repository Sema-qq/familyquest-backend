package user

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

//go:embed sqls/get_by_login.sql
var getByLoginSQL string

//go:embed sqls/get_by_id.sql
var getByIDSQL string

type Repository struct {
	db db.Conn
}

func NewRepository(db db.Conn) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, user entity.UserCreate) error {
	_, err := r.db.Exec(
		ctx,
		createSQL,
		user.ID.UUID(),
		user.Login,
		user.PasswordHash,
		user.DisplayName,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.AlreadyExists("user already exists")
		}

		return fmt.Errorf("failed create user: %w", err)
	}

	return nil
}

func (r *Repository) GetByLogin(ctx context.Context, login string) (entity.UserCredentials, error) {
	user, err := db.QueryRow[userCredentials](ctx, r.db, getByLoginSQL, login)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.UserCredentials{}, domain.NotFound("user not found")
		}

		return entity.UserCredentials{}, fmt.Errorf("failed get user by login: %w", err)
	}

	return user.toEntity(), nil
}

func (r *Repository) GetByID(ctx context.Context, id entity.UserID) (entity.User, error) {
	user, err := db.QueryRow[user](ctx, r.db, getByIDSQL, id.UUID())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, domain.NotFound("user not found")
		}

		return entity.User{}, fmt.Errorf("failed get user by id: %w", err)
	}

	return user.toEntity(), nil
}
