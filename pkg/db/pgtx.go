package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type pgTx struct {
	tx pgx.Tx
}

func newPgTx(tx pgx.Tx) Conn {
	return &pgTx{
		tx: tx,
	}
}

func (p *pgTx) Ping(context.Context) error {
	return nil
}

func (p *pgTx) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return p.tx.Exec(ctx, sql, args...)
}

func (p *pgTx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return p.tx.Query(ctx, sql, args...)
}

func (p *pgTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return p.tx.QueryRow(ctx, sql, args...)
}

func (p *pgTx) BeginTxFunc(ctx context.Context, _ pgx.TxOptions, f func(Conn) error) error {
	return pgx.BeginFunc(ctx, p.tx, func(tx pgx.Tx) error {
		return f(newPgTx(tx))
	})
}

func (p *pgTx) Close() {}
