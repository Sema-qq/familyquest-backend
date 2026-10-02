package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgPool struct {
	pool *pgxpool.Pool
}

func newPgPool(ctx context.Context, config *pgxpool.Config) (Conn, error) {
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &pgPool{
		pool: pool,
	}, nil
}

func (p *pgPool) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

func (p *pgPool) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return p.pool.Exec(ctx, sql, args...)
}

func (p *pgPool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return p.pool.Query(ctx, sql, args...)
}

func (p *pgPool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return p.pool.QueryRow(ctx, sql, args...)
}

func (p *pgPool) BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, f func(Conn) error) error {
	return pgx.BeginTxFunc(ctx, p.pool, txOptions, func(tx pgx.Tx) error {
		return f(newPgTx(tx))
	})
}

func (p *pgPool) Close() {
	p.pool.Close()
}
