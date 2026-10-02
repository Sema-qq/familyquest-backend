package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type txAwareConn struct {
	base Conn
}

func NewTxAwareConn(base Conn) Conn {
	return &txAwareConn{
		base: base,
	}
}

func (c *txAwareConn) conn(ctx context.Context) Conn {
	if tx, ok := txFromContext(ctx); ok {
		return tx
	}

	return c.base
}

func (c *txAwareConn) Ping(ctx context.Context) error {
	return c.base.Ping(ctx)
}

func (c *txAwareConn) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return c.conn(ctx).Exec(ctx, sql, args...)
}

func (c *txAwareConn) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return c.conn(ctx).Query(ctx, sql, args...)
}

func (c *txAwareConn) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return c.conn(ctx).QueryRow(ctx, sql, args...)
}

func (c *txAwareConn) BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, f func(Conn) error) error {
	return c.conn(ctx).BeginTxFunc(ctx, txOptions, f)
}

func (c *txAwareConn) Close() {
	c.base.Close()
}
