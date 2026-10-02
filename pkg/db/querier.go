package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Querier interface {
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
}

func QueryRow[T any](ctx context.Context, db Querier, query string, args ...any) (T, error) {
	var record T

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return record, fmt.Errorf("query: %w", err)
	}

	record, err = pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[T])
	if err != nil {
		return record, fmt.Errorf("scan row: %w", err)
	}

	return record, nil
}

func Query[T any](ctx context.Context, db Querier, query string, args ...any) ([]T, error) {
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	records, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[T])
	if err != nil {
		return nil, fmt.Errorf("scan rows: %w", err)
	}

	return records, nil
}
