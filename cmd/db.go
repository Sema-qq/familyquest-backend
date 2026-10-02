package main

import (
	"context"

	"familyquest-backend/pkg/db"
)

func NewPostgresConn(ctx context.Context, cfg Config) (db.Conn, error) {
	conn, err := db.NewPostgresConn(ctx, NewPostgresConfig(cfg))
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(ctx); err != nil {
		conn.Close()
		return nil, err
	}

	return db.NewTxAwareConn(conn), nil
}
