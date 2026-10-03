package familymember

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

//go:embed sqls/get_by_user.sql
var getByUserSQL string

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

func (r *Repository) Create(ctx context.Context, member entity.FamilyMemberCreate) error {
	_, err := r.db.Exec(
		ctx,
		createSQL,
		member.ID.UUID(),
		member.FamilyID.UUID(),
		member.UserID.UUID(),
		string(member.Role),
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.Conflict("family member already exists")
		}

		return fmt.Errorf("failed create family member: %w", err)
	}

	return nil
}

func (r *Repository) GetByUser(ctx context.Context, userID entity.UserID) (entity.FamilyMember, error) {
	member, err := db.QueryRow[familyMember](ctx, r.db, getByUserSQL, userID.UUID())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.FamilyMember{}, domain.NotFound("family member not found")
		}

		return entity.FamilyMember{}, fmt.Errorf("failed get family member by user: %w", err)
	}

	return member.toEntity(), nil
}

func (r *Repository) ListByFamily(ctx context.Context, familyID entity.FamilyID) ([]entity.FamilyMemberProfile, error) {
	members, err := db.Query[familyMemberProfile](ctx, r.db, listByFamilySQL, familyID.UUID())
	if err != nil {
		return nil, fmt.Errorf("failed list family members by family: %w", err)
	}

	return familyMemberProfileList(members).toEntities(), nil
}
