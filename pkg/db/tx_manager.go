package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type TxManager struct {
	db Conn
}

func NewTxManager(db Conn) *TxManager {
	return &TxManager{
		db: db,
	}
}

func (m *TxManager) WithinTx(ctx context.Context, fn func(context.Context) error) error {
	return m.WithinTxOptions(ctx, pgx.TxOptions{}, fn)
}

func (m *TxManager) WithinTxOptions(
	ctx context.Context,
	txOptions pgx.TxOptions,
	fn func(context.Context) error,
) error {
	if hasTx(ctx) {
		return fn(ctx)
	}

	return m.db.BeginTxFunc(ctx, txOptions, func(tx Conn) error {
		return fn(withTx(ctx, tx))
	})
}
