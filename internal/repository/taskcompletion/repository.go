package taskcompletion

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

//go:embed sqls/get_rejected_by_member_and_slot.sql
var getRejectedByMemberAndSlotSQL string

//go:embed sqls/resubmit_rejected.sql
var resubmitRejectedSQL string

//go:embed sqls/list_by_season.sql
var listBySeasonSQL string

//go:embed sqls/get_review_by_id.sql
var getReviewByIDSQL string

//go:embed sqls/approve.sql
var approveSQL string

//go:embed sqls/reject.sql
var rejectSQL string

//go:embed sqls/count_submitted_by_season.sql
var countSubmittedBySeasonSQL string

//go:embed sqls/results_by_season.sql
var resultsBySeasonSQL string

//go:embed sqls/has_approved_by_slot.sql
var hasApprovedBySlotSQL string

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

func (r *Repository) Create(ctx context.Context, completion entity.TaskCompletionCreate) (entity.TaskCompletion, error) {
	result, err := db.QueryRow[taskCompletion](
		ctx,
		r.db,
		createSQL,
		completion.ID.UUID(),
		completion.MemberID.UUID(),
		completion.SlotID.UUID(),
		completion.PerformedOn,
		completion.ChildComment,
	)
	if err != nil {
		return entity.TaskCompletion{}, fmt.Errorf("failed create task completion: %w", err)
	}

	return result.toEntity(), nil
}

func (r *Repository) GetRejectedByMemberAndSlot(
	ctx context.Context,
	memberID entity.FamilyMemberID,
	slotID entity.SeasonTaskSlotID,
) (entity.TaskCompletion, bool, error) {
	result, err := db.QueryRow[taskCompletion](
		ctx,
		r.db,
		getRejectedByMemberAndSlotSQL,
		memberID.UUID(),
		slotID.UUID(),
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.TaskCompletion{}, false, nil
		}

		return entity.TaskCompletion{}, false, fmt.Errorf("failed get rejected task completion by member and slot: %w", err)
	}

	return result.toEntity(), true, nil
}

func (r *Repository) ResubmitRejected(
	ctx context.Context,
	completion entity.TaskCompletionResubmit,
) (entity.TaskCompletion, error) {
	result, err := db.QueryRow[taskCompletion](
		ctx,
		r.db,
		resubmitRejectedSQL,
		completion.ID.UUID(),
		completion.PerformedOn,
		completion.ChildComment,
	)
	if err != nil {
		return entity.TaskCompletion{}, fmt.Errorf("failed resubmit rejected task completion: %w", err)
	}

	return result.toEntity(), nil
}

func (r *Repository) ListBySeason(
	ctx context.Context,
	seasonID entity.SeasonID,
	memberID *entity.FamilyMemberID,
	status *entity.TaskCompletionStatus,
) ([]entity.TaskCompletionView, error) {
	var memberArg any
	if memberID != nil {
		memberArg = memberID.UUID()
	}
	var statusArg any
	if status != nil {
		statusArg = string(*status)
	}

	result, err := db.Query[taskCompletionView](ctx, r.db, listBySeasonSQL, seasonID.UUID(), memberArg, statusArg)
	if err != nil {
		return nil, fmt.Errorf("failed list task completions by season: %w", err)
	}

	return taskCompletionViewList(result).toEntities(), nil
}

func (r *Repository) GetReviewByID(ctx context.Context, id entity.TaskCompletionID) (entity.TaskCompletionReview, error) {
	result, err := db.QueryRow[taskCompletionReview](ctx, r.db, getReviewByIDSQL, id.UUID())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.TaskCompletionReview{}, domain.NotFound("task completion not found")
		}

		return entity.TaskCompletionReview{}, fmt.Errorf("failed get task completion review by id: %w", err)
	}

	return result.toEntity(), nil
}

func (r *Repository) Approve(
	ctx context.Context,
	update entity.TaskCompletionUpdateReview,
) (entity.TaskCompletion, error) {
	result, err := db.QueryRow[taskCompletion](
		ctx,
		r.db,
		approveSQL,
		update.ID.UUID(),
		update.ReviewedBy.UUID(),
		update.ReviewedAt,
	)
	if err != nil {
		return entity.TaskCompletion{}, fmt.Errorf("failed approve task completion: %w", err)
	}

	return result.toEntity(), nil
}

func (r *Repository) Reject(
	ctx context.Context,
	update entity.TaskCompletionUpdateReview,
) (entity.TaskCompletion, error) {
	result, err := db.QueryRow[taskCompletion](
		ctx,
		r.db,
		rejectSQL,
		update.ID.UUID(),
		update.ParentComment,
		update.ReviewedBy.UUID(),
		update.ReviewedAt,
	)
	if err != nil {
		return entity.TaskCompletion{}, fmt.Errorf("failed reject task completion: %w", err)
	}

	return result.toEntity(), nil
}

func (r *Repository) CountSubmittedBySeason(ctx context.Context, seasonID entity.SeasonID) (int64, error) {
	result, err := db.QueryRow[count](ctx, r.db, countSubmittedBySeasonSQL, seasonID.UUID())
	if err != nil {
		return 0, fmt.Errorf("failed count submitted task completions by season: %w", err)
	}

	return result.Count, nil
}

func (r *Repository) ResultsBySeason(
	ctx context.Context,
	seasonID entity.SeasonID,
	memberID *entity.FamilyMemberID,
) ([]entity.SeasonResult, error) {
	var memberArg any
	if memberID != nil {
		memberArg = memberID.UUID()
	}

	result, err := db.Query[seasonResult](ctx, r.db, resultsBySeasonSQL, seasonID.UUID(), memberArg)
	if err != nil {
		return nil, fmt.Errorf("failed get season results: %w", err)
	}

	return seasonResultList(result).toEntities(), nil
}

func (r *Repository) HasApprovedBySlot(ctx context.Context, slotID entity.SeasonTaskSlotID) (bool, error) {
	result, err := db.QueryRow[count](
		ctx,
		r.db,
		hasApprovedBySlotSQL,
		slotID.UUID(),
	)
	if err != nil {
		return false, fmt.Errorf("failed check approved task completion by slot: %w", err)
	}

	return result.Count > 0, nil
}
