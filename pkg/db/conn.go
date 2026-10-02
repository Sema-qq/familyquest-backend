package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Conn interface {
	ReadConn
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, f func(Conn) error) error
}

type ReadConn interface {
	Ping(ctx context.Context) error
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Close()
}
